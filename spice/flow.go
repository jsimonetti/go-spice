package spice

import (
	"io"
	"net"
	"sync"
	"time"
)

type Flow struct {
	tenant  net.Conn
	compute net.Conn
}

func NewFlow(tenant net.Conn, compute net.Conn) *Flow {

	flow := &Flow{
		tenant:  tenant,
		compute: compute,
	}

	return flow
}

func (f *Flow) proxyData(src net.Conn, dst net.Conn) error {
	data := make([]byte, 64*1024)
	pending := 0
	for {
		if pending == 0 {
			var err error
			pending, err = src.Read(data)
			if err != nil && err != io.EOF {
				return err
			}
			if pending == 0 {
				return nil
			}
		}

		done, err := dst.Write(data[0:pending])
		if err != nil {
			return err
		}
		data = data[done:]
		pending -= done
	}
}

func (f *Flow) proxyToCompute() error {
	err := f.proxyData(f.tenant, f.compute)
	f.compute.Close()
	return err
}

func (f *Flow) proxyToTenant() error {
	err := f.proxyData(f.compute, f.tenant)
	f.tenant.Close()
	return err
}

func (f *Flow) Proxy() error {
	Pipe(f.compute, f.tenant)
	return nil

	go f.proxyToTenant()
	return f.proxyToCompute()
}

func writePending(data []byte, dst io.Writer) error {
	length := len(data)
	for length > 1 {
		done, err := dst.Write(data[0:length])
		if err != nil {
			return err
		}
		data = data[done:]
		length -= done
	}
	return nil
}

func Pipe(src, dst net.Conn) (sent, received int64) {
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		sent = PipeAndClose(src, dst)
		wg.Done()
	}()
	go func() {
		received = PipeAndClose(dst, src)
		wg.Done()
	}()
	wg.Wait()
	return
}

func PipeAndClose(src, dst net.Conn) (copied int64) {
	timeout := 10 * time.Second
	defer dst.Close()
	buf := make([]byte, 4096)
	for {
		n, err := src.Read(buf)
		copied += int64(n)
		if n > 0 {
			dst.SetWriteDeadline(time.Now().Add(timeout))
			if _, err := dst.Write(buf[0:n]); err != nil {
				break
			}
		}
		if err != nil {
			break
		}
	}
	return
}

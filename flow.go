package spice

import (
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

func (f *Flow) Proxy() error {
	f.pipe(f.compute, f.tenant)
	return nil
}

func (f *Flow) pipe(src, dst net.Conn) (sent, received int64) {
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		sent = f.pipeAndClose(src, dst)
		wg.Done()
	}()
	go func() {
		received = f.pipeAndClose(dst, src)
		wg.Done()
	}()
	wg.Wait()
	return
}

func (f *Flow) pipeAndClose(src, dst net.Conn) (copied int64) {
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

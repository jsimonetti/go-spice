package spice

import (
	"net"
	"sync"

	"acln.ro/zerocopy"
)

// flow is a connection pipe to couple tenant to compute connections
type flow struct {
	tenant  net.Conn
	compute net.Conn
}

// newFlow returns a new flow
func newFlow(tenant net.Conn, compute net.Conn) *flow {
	flow := &flow{
		tenant:  tenant,
		compute: compute,
	}
	return flow
}

// Pipe will start piping the connections together
func (f *flow) Pipe() error {
	f.pipe(f.compute, f.tenant)
	return nil
}

func (f *flow) pipe(src, dst net.Conn) (sent, received int64) {
	if src == nil || dst == nil {
		return
	}

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		sent, _ = zerocopy.Transfer(src, dst)
		wg.Done()
	}()
	go func() {
		received, _ = zerocopy.Transfer(dst, src)
		wg.Done()
	}()
	wg.Wait()
	return
}

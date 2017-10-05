go-spice [![GoDoc](https://godoc.org/github.com/jsimonetti/go-spice?status.svg)](https://godoc.org/github.com/jsimonetti/go-spice) [![Go Report Card](https://goreportcard.com/badge/github.com/jsimonetti/go-spice)](https://goreportcard.com/report/github.com/jsimonetti/go-spice)
=======

Package `spice` attempts to implement a SPICE proxy.
It can be used to proxy virt-viewer/remote-viewer traffic to destination qemu instances.

This package is still unfinished. The API is highly unstable.

TODO:
- implement proper auth capability handling
- update documentation

Not planned, but nice to have
- implement SASL authentication


See [example](examples/proxy.go) for an example including an Authenticator
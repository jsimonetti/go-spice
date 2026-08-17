module github.com/jsimonetti/go-spice

go 1.23

replace acln.ro/zerocopy => github.com/acln0/zerocopy v0.0.0-20190410132315-ac749309e897

require (
	acln.ro/zerocopy v0.0.0-20190410132315-ac749309e897
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.10.0
)

require golang.org/x/sys v0.13.0 // indirect

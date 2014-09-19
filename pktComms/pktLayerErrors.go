package pktComms

import (
	e "errors"
)

var (
	NilCnx              = e.New("nil tcp connection")
	NilNode             = e.New("nil Node for PktLayer")
	UnusableInCnx       = e.New("nil or otherwise unusable in connection")
	UnusableTcpAcceptor = e.New("nil or otherwise unusable tcp acceptor")
)

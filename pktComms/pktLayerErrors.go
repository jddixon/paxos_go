package pktComms

import (
	e "errors"
)

var (
	NilNode       = e.New("nil Node for PktLayer")
	UnusableInCnx = e.New("nil or otherwise unusable in connection")
)

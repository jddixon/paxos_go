package pktComms

import (
	. "gopkg.in/check.v1"
	"testing"
	"time"
)

// IF USING gocheck, need a file like this in each package=directory.

func Test(t *testing.T) { TestingT(t) }

type XLSuite struct{}

var _ = Suite(&XLSuite{})

const (
	KA_COUNT    = 15 // number of keepalives to send in tests
	KA_INTERVAL = 200 * time.Millisecond
)

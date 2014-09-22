package pktComms

// paxos_go/pktComms/outQ.go

import (
	xt "github.com/jddixon/xlTransport_go"
	"sync"
	"time"
)

type OutMsg struct {
	Packet  []byte
	Timeout time.Duration
	Rcvd    time.Time
	Sender  chan *OutResponse // where to forward the reply
}

// Reply to an OutMsg
type OutResponse struct {
	Data     []byte
	Timedout time.Duration // if zero, it didn't time out
	Rcvd     time.Time
}

type OutQ struct {
	Msgs []*OutMsg
	mu   sync.RWMutex
}

type OutCnxMgr struct {
	Cnx    *xt.TcpConnection
	MsqQ   *OutQ
	CurMsg *OutMsg
	mu     sync.RWMutex
	State  uint
}

// CnxTalker a useful abstraction ?  CnxTalkerFactory?

package pktComms

// paxos_go/pktComms/inQ.go

import (
	xn "github.com/jddixon/xlNode_go"
	xt "github.com/jddixon/xlTransport_go"
	"sync"
	"time"
)

type InMsg struct {
	Packet []byte
	Rcvd   time.Time
}

// Reply to an InMsg
type InResponse struct {
	Data []byte
}

type InQ struct {
	Msgs []*InMsg
	mu   sync.RWMutex
}

func NewInQ() (inQ *InQ, err error) {
	inQ = &InQ{}
	return
}

type InCnxMgr struct {
	Cnx        *xt.TcpConnection
	Introduced bool // true if Hello rcvd but no Bye yet
	MsgQ       *InQ
	CurMsg     *InMsg
	mu         sync.RWMutex
	State      uint
	MyNode     *xn.Node
}

func NewInCnxMgr(node *xn.Node, cnx *xt.TcpConnection) (
	icm *InCnxMgr, err error) {

	var inQ *InQ

	if node == nil {
		err = NilNode
	} else if cnx == nil {
		err = UnusableInCnx
	} else {
		inQ, err = NewInQ()
		if err == nil {
			icm = &InCnxMgr{
				Cnx:    cnx,
				MsgQ:   inQ,
				MyNode: node,
			}
		}
	}
	return
}

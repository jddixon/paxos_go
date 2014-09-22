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

type InListener struct {
	Acc    *xt.TcpAcceptor
	MyNode *xn.Node
}

// XXX This should have an InCnxMgr = CnxListener factory as an argument.
func NewInListener(node *xn.Node, acc *xt.TcpAcceptor) (
	inL *InListener, err error) {

	if node == nil {
		err = NilNode
	} else if acc == nil {
		err = UnusableTcpAcceptor
	} else {
		inL = &InListener{
			Acc:    acc,
			MyNode: node,
		}
	}
	return
}

// XXX A mistake ...
func (inL *InListener) Accept() (icm *InCnxMgr, err error) {
	conn, err := inL.Acc.Accept()
	if err == nil {
		if conn == nil {
			err = NilCnx
		}
		cnx := conn.(*xt.TcpConnection)
		icm, err = NewInCnxMgr(inL.MyNode, cnx)
	}
	return
}

func (inL *InListener) Run() {

	// Listen on Acc.  If someone opens a connection, create an inQMgr
	// and run it in a separate goroutine.

	// We could keep a slice of connections opened.  This would not be a
	// good idea if there were to be many such connections.

}

func (inL *InListener) Close() (err error) {
	if inL.Acc != nil {
		inL.Acc.Close()
	}
	return // XXX err not changed
}

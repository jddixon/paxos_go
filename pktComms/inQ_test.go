package pktComms

// paxos_go/pktComms/inQ_test.go

import (
	"fmt"
	xr "github.com/jddixon/rnglib_go"
	xi "github.com/jddixon/xlNodeID_go"
	xn "github.com/jddixon/xlNode_go"
	xt "github.com/jddixon/xlTransport_go"
	. "gopkg.in/check.v1"
)

func (s *XLSuite) dummyPeer(c *C, acc *xt.TcpAcceptor) {

	// open a connection to the listening socket, then send it some
	// number (15?) of KeepAlive packets, expecting an Ack for each

}
func (s *XLSuite) TestInQ(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("\nTEST_IN_Q")
	}

	rng := xr.MakeSimpleRNG()

	inQ, err := NewInQ()
	c.Assert(err, IsNil)
	c.Assert(len(inQ.Msgs), Equals, 0)

	nodeName := rng.NextFileName(8)
	lfs, err := s.makeUniqueLFS(c, rng)
	c.Assert(err, IsNil)
	c.Assert(lfs, NotNil)

	id, err := xi.New(nil)
	c.Assert(err, IsNil)
	c.Assert(id, NotNil)

	node, err := xn.NewNew(nodeName, id, lfs)
	c.Assert(err, IsNil)
	c.Assert(node, NotNil)

	// create an acceptor, a listening socket, on a random localhost port
	acc, err := xt.NewTcpAcceptor("127.0.0.1:0")
	c.Assert(err, IsNil)
	c.Assert(acc, NotNil)
	defer acc.Close()

	inL, err := NewInListener(node, acc)
	c.Assert(err, IsNil)
	c.Assert(inL, NotNil)
	c.Assert(inL.Acc, DeepEquals, acc)
	c.Assert(inL.MyNode, DeepEquals, node)
}
func (s *XLSuite) TestInMain(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("\nTEST_IN_MAIN")
	}

	rng := xr.MakeSimpleRNG()

	nodeName := rng.NextFileName(8)
	lfs, err := s.makeUniqueLFS(c, rng)
	c.Assert(err, IsNil)
	c.Assert(lfs, NotNil)

	// DEBUG
	fmt.Printf("lfs = %s\n", lfs)
	// END
	id, err := xi.New(nil)
	c.Assert(err, IsNil)
	c.Assert(id, NotNil)

	node, err := xn.NewNew(nodeName, id, lfs)
	c.Assert(err, IsNil)
	c.Assert(node, NotNil)

	// create an acceptor, a listening socket, on a random localhost port
	acc, err := xt.NewTcpAcceptor("127.0.0.1:0")
	c.Assert(err, IsNil)
	c.Assert(acc, NotNil)
	defer acc.Close()

	inL, err := NewInListener(node, acc)
	c.Assert(err, IsNil)
	c.Assert(inL, NotNil)

	// Create a cluster of K nodes for test purposes.

	// Run the listener in a separate goroutine.  This will spawn an
	// inQMgr for each connection.  Each inQMgr simulates a peer in the
	// cluster; we don't simulate peer SelfIndex

	// go inL.Run()

	// Opens K-1 connections to listener inL.
	// Each such connection is passed to an InQDriver(cnx) which does
	//   Hello/Ack
	//   N * KeepAlive/Ack, where X = say 15
	//   Bye/Ack

	// The node in InListener is the SelfIndex node.

	// There is then an InQDriver for each of the other K-1 nodes in
	// the cluster.

}

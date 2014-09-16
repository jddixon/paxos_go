package pktComms

// paxos_go/pktComms/helloGoodbye_test.go

import (
	"fmt"
	xr "github.com/jddixon/rnglib_go"
	//xg "github.com/jddixon/xlReg_go"
	. "gopkg.in/check.v1"
)

// Launch K servers, have them say hello to one another, pause, then
// have them say goodbye.
func (s *XLSuite) TestHelloGoodbye(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("TEST_HELLO_GOODBYE")
	}
	rng := xr.MakeSimpleRNG()

	// 1. Launch an ephemeral xlReg server --------------------------
	eph, reg, regID, server := s.launchEphServer(c)
	defer eph.Close()

	// 2. Create a random cluster name and size; register it --------
	clusterName, clusterAttrs, clusterID, K := s.createAndRegSoloCluster(
		c, rng, reg, regID, server)

	// 3  Create K cluster member PktLayers
	pL, pLNames := s.createKMemberPktLayers(c, rng, server,
		clusterName, clusterAttrs, clusterID, K)

	_, _ = pL, pLNames // XXX not yet used

	// XXX This is a hack: we should wait until all have joined, then
	// wait till all have said hello, then wait until all have said
	// Bye.  There is a DoneCh in MemberNode.

	// 4  Start the K clients running, each in a separate goroutine
	for i := uint32(0); i < K; i++ {
		pL[i].Run()
	}

}

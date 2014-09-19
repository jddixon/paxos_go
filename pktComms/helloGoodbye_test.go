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

	// 4  Start the K clients running, each in a separate goroutine.
	for i := uint32(0); i < K; i++ {
		go pL[i].JoinCluster()
	}
	for i := uint32(0); i < K; i++ {
		ok := <-pL[i].DoneCh
		// DEBUG
		fmt.Printf("member %d, %-8s,  has joined ", i, pLNames[i])
		if ok {
			fmt.Println("successfully")
		} else {
			// XXX Using pL.Err will cause timing problems
			fmt.Printf("but returned an error %s\n", pL[i].Err)
		}
		// END
	}

	// 5  Tell all to say Hello; wait.

	// 6  Tell all to say Byte; wait.

	// 7  We are done.

}

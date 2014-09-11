package pktComms

// paxos_go/pktComms/keep_alive_test.go

import (
	"fmt"
	xr "github.com/jddixon/rnglib_go"
	xg "github.com/jddixon/xlReg_go"
	. "gopkg.in/check.v1"
)

/////////////////////////////////////////////////////////////////////
// XXX THIS IS BEING HACKED INTO A TEST OF A CLUSTER JUST RUNNING pktComms
// KEEP-ALIVE MESSAGES
/////////////////////////////////////////////////////////////////////

func (s *XLSuite) TestKeepAlives(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("TEST_KEEP_ALIVES")
	}
	rng := xr.MakeSimpleRNG()

	/////////////////////////////////////////////////////////////////
	// A: Launch N tcNodes for cluster cl to coordinate through
	// xlReg at 127.0.0.1:PPPPP
	/////////////////////////////////////////////////////////////////

	// we listen on three ports: command, intra-cluster comms, and
	// a third for external clients
	epCount := uint(3)
	maxSize := uint(2 + rng.Intn(6)) // so from 2 to 7
	cl := s.makeACluster(c, rng, epCount, maxSize)

	c.Assert(cl.MaxSize(), Equals, maxSize)
	c.Assert(cl.Size(), Equals, maxSize) //

	// Verify that member names are unique within the cluster
	ids := make([][]byte, maxSize)
	names := make([]string, maxSize)
	nameMap := make(map[string]uint)
	for i := uint(0); i < maxSize; i++ {
		member := cl.Members[i]
		names[i] = member.GetName()
		nameMap[names[i]] = i

		// collect IDs while we are at it
		id := member.GetNodeID().Value() // returns a clone of the nodeID
		ids[i] = id
	}
	// if the names are not unique, map will be smaller
	c.Assert(maxSize, Equals, uint(len(nameMap)))

	// verify that the RegCluster.MembersByName index is correct
	for i := uint(0); i < maxSize; i++ {
		name := names[i]
		member := cl.MembersByName[name]
		c.Assert(name, Equals, member.GetName())
	}

	// verify that the RegCluster.MembersByID index is correct
	count := uint(0) // number of successful type assertions
	for i := uint(0); i < maxSize; i++ {
		id := ids[i]
		mbr, err := cl.MembersByID.Find(id)
		c.Assert(err, IsNil)
		var member *xg.MemberInfo
		// verify that the type assertion succeeds
		if m, ok := mbr.(*xg.MemberInfo); ok {
			member = m
			mID := member.GetNodeID().Value()
			c.Assert(len(id), Equals, len(mID))
			for j := uint(0); j < uint(len(id)); j++ {
				c.Assert(id[j], Equals, mID[j])
			}
			count++
		}
	}
	c.Assert(maxSize, Equals, count)

	/////////////////////////////////////////////////////////////////
	// B: Each tcNode confirgures acceptor An = a random tcpip
	// endpoint 127.0.0.1:Pn; selects keys sPriv, cPriv; selects
	// nodeID = SHA3(127.0.0.1:Pn (+) sPriv (+) cPriv
	/////////////////////////////////////////////////////////////////

	// XXX STUB

	/////////////////////////////////////////////////////////////////
	// C: Each tcNode initiates xlReg cycle, at end of which N-1 peers
	// are configured.
	/////////////////////////////////////////////////////////////////

	// XXX STUB

	/////////////////////////////////////////////////////////////////
	// D: Each tcNode starts N-1 Hello/Ack cycles with peers.
	/////////////////////////////////////////////////////////////////

	// XXX STUB

	/////////////////////////////////////////////////////////////////
	// E: When all N-1 Hellos have been Acked, each tcNode initiates
	// K = 15 KeepAlive/Ack cycles with its N-1 peers.  Pause 2 sec
	// between an Ack and the next KeepAlive.
	/////////////////////////////////////////////////////////////////

	// XXX STUB

	/////////////////////////////////////////////////////////////////
	// F: When K=15 KeepAlive/Ack cycles have been completed with a
	// peer, each tcNode waits 2 seconds and then sends a Bye to the
	// peer.  When an Ack to the Bye has been received, the tcNode
	// closes that connection.
	/////////////////////////////////////////////////////////////////

	// XXX STUB

	/////////////////////////////////////////////////////////////////
	// G: When N-1 Bye/Ack cycles have been completed, tcNode sends
	// stopped to the controller.
	/////////////////////////////////////////////////////////////////

	// XXX STUB

	/////////////////////////////////////////////////////////////////
	// H: When the controller has received N stopped signals, the
	// test is over.
	/////////////////////////////////////////////////////////////////

	// XXX STUB

}
func (s *XLSuite) TestClusterSerialization(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("TEST_CLUSTER_SERIALIZATION")
	}
	rng := xr.MakeSimpleRNG()

	// Generate a random cluster
	epCount := uint(1 + rng.Intn(3)) // so from 1 to 3
	size := uint(2 + rng.Intn(6))    // so from 2 to 7
	cl := s.makeACluster(c, rng, epCount, size)

	// Serialize it
	serialized := cl.String()

	// Reverse the serialization
	deserialized, rest, err := xg.ParseRegCluster(serialized)
	c.Assert(err, IsNil)
	c.Assert(deserialized, Not(IsNil))
	c.Assert(len(rest), Equals, 0)

	// Verify that the deserialized cluster is identical to the original
	c.Assert(deserialized.Equal(cl), Equals, true)

}

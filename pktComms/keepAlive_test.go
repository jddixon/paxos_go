package pktComms

// paxos_go/pktComms/keep_alive_test.go

import (
	"fmt"
	xr "github.com/jddixon/rnglib_go"
	//xg "github.com/jddixon/xlReg_go"
	. "gopkg.in/check.v1"
)

func (s *XLSuite) TestKeepAlives(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("\nTEST_KEEP_ALIVES")
	}
	rng := xr.MakeSimpleRNG()

	// 1. Launch an ephemeral xlReg server --------------------------
	eph, reg, regID, server := s.launchEphServer(c)
	defer eph.Close()

	// 2. Create a random cluster name and size ---------------------
	clusterName, clusterAttrs, clusterID, K := s.createAndRegSoloCluster(c, rng,
		reg, regID, server)

	_, _, _, _ = clusterName, clusterAttrs, clusterID, K // XXX NOT YET USED

	////////////////////////////////////////////////////////////////////
	// XXX COMPLETE helloGoodbye_test.go, THEN MODIFY WHAT FOLLOWS IN
	// SAME SPIRIT
	////////////////////////////////////////////////////////////////////

	/////////////////////////////////////////////////////////////////
	// B: Launch N tcNodes for cluster cl to coordinate through
	// xlReg at 127.0.0.1:PPPPP.   Each tcNode configures acceptor
	// An = a random tcpip endpoint 127.0.0.1:Pn; selects keys sPriv, cPriv
	/////////////////////////////////////////////////////////////////

	//// we listen on three ports: command, intra-cluster comms, and
	//// a third for external clients
	//epCount := uint32(3)
	//maxSize := uint32(2 + rng.Intn(6)) // so from 2 to 7
	//cl, nodes, ckPrivs, skPrivs := s.makeACluster(c, rng, epCount, maxSize)

	//// XXX nodes, key slices not currently used
	//_, _, _ = nodes, ckPrivs, skPrivs

	//c.Assert(cl.MaxSize(), Equals, maxSize)
	//c.Assert(cl.Size(), Equals, maxSize)
	//c.Assert(maxSize, Equals, uint32(len(nodes)))

	//// Verify that member names are unique within the cluster
	//ids := make([][]byte, maxSize)
	//names := make([]string, maxSize)
	//nameMap := make(map[string]uint32)
	//for i := uint32(0); i < maxSize; i++ {
	//	member := cl.Members[i]
	//	names[i] = member.GetName()
	//	nameMap[names[i]] = i

	//	// collect IDs while we are at it
	//	id := member.GetNodeID().Value() // returns a clone of the nodeID
	//	ids[i] = id
	//}
	//// if the names are not unique, map will be smaller
	//c.Assert(maxSize, Equals, uint32(len(nameMap)))

	//// verify that the RegCluster.MembersByName index is correct
	//for i := uint32(0); i < maxSize; i++ {
	//	name := names[i]
	//	member := cl.MembersByName[name]
	//	c.Assert(name, Equals, member.GetName())
	//}

	//// verify that the RegCluster.MembersByID index is correct
	//count := uint32(0) // number of successful type assertions
	//for i := uint32(0); i < maxSize; i++ {
	//	id := ids[i]
	//	mbr, err := cl.MembersByID.Find(id)
	//	c.Assert(err, IsNil)
	//	var member *xg.MemberInfo
	//	// verify that the type assertion succeeds
	//	if m, ok := mbr.(*xg.MemberInfo); ok {
	//		member = m
	//		mID := member.GetNodeID().Value()
	//		c.Assert(len(id), Equals, len(mID))
	//		for j := uint(0); j < uint(len(id)); j++ {
	//			c.Assert(id[j], Equals, mID[j])
	//		}
	//		count++
	//	}
	//}
	//c.Assert(maxSize, Equals, count)

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

//func (s *XLSuite) TestClusterSerialization(c *C) {
//	if VERBOSITY > 0 {
//		fmt.Println("\nTEST_CLUSTER_SERIALIZATION")
//	}
//	rng := xr.MakeSimpleRNG()
//
//	// Generate a random cluster
//	epCount := uint32(1 + rng.Intn(3)) // so from 1 to 3
//	size := uint32(2 + rng.Intn(6))    // so from 2 to 7
//	// XXX MEMBERS, KEY SLICES NOT YET USED
//	cl, _, _, _ := s.makeACluster(c, rng, epCount, size)
//
//	// Serialize it
//	serialized := cl.String()
//
//	// Reverse the serialization
//	deserialized, rest, err := xg.ParseRegCluster(serialized)
//	c.Assert(err, IsNil)
//	c.Assert(deserialized, Not(IsNil))
//	c.Assert(len(rest), Equals, 0)
//
//	// Verify that the deserialized cluster is identical to the original
//	c.Assert(deserialized.Equal(cl), Equals, true)
//
//}

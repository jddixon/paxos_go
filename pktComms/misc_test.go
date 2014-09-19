package pktComms

// paxos_go/pktComms/misc_test.go

import (
	//"crypto/rand"
	//"crypto/rsa"
	//"encoding/hex"
	"fmt"
	xr "github.com/jddixon/rnglib_go"
	xi "github.com/jddixon/xlNodeID_go"
	//xn "github.com/jddixon/xlNode_go"
	xg "github.com/jddixon/xlReg_go"
	xt "github.com/jddixon/xlTransport_go"
	. "gopkg.in/check.v1"
	//"strings"
)

const (
	VERBOSITY = 1
)

// Create and start an ephemeral xlReg server -----------------------
// Calling routine should call defer eph.Close()
func (s *XLSuite) launchEphServer(c *C) (eph *xg.EphServer,
	reg *xg.Registry, regID *xi.NodeID, server *xg.RegServer) {

	eph, err := xg.NewEphServer()
	c.Assert(eph, NotNil)
	c.Assert(err, IsNil)

	server = eph.Server

	// start the ephemeral server -------------------------
	err = eph.Run()
	c.Assert(err, IsNil)

	// verify Bloom filter is running
	reg = &eph.Server.Registry
	c.Assert(reg, NotNil)
	regID = reg.GetNodeID()
	c.Assert(reg.IDCount(), Equals, uint(1)) // the registry's own ID
	found, err := reg.ContainsID(regID)
	c.Assert(err, IsNil)
	c.Assert(found, Equals, true)

	return
}

// Create and register a solo cluster -------------------------------
func (s *XLSuite) createAndRegSoloCluster(c *C, rng *xr.PRNG,
	reg *xg.Registry, regID *xi.NodeID, server *xg.RegServer) (
	clusterName string, clusterAttrs uint64, clusterID *xi.NodeID, K uint32) {

	serverName := server.GetName()
	serverID := server.GetNodeID()
	serverEnd := server.GetEndPoint(0)
	serverCK := server.GetCommsPublicKey()
	serverSK := server.GetSigPublicKey()
	c.Assert(serverEnd, NotNil)

	clusterName = rng.NextFileName(8)
	clusterAttrs = uint64(rng.Int63())
	K = uint32(2 + rng.Intn(6)) // so the size is 2 .. 7

	// create an AdminMember, use it to get the clusterID --------
	an, err := xg.NewAdminMember(serverName, serverID, serverEnd,
		serverCK, serverSK, clusterName, clusterAttrs, K, uint32(3), nil)
	c.Assert(err, IsNil)

	an.Run()
	<-an.DoneCh

	c.Assert(an.ClusterID, NotNil)          // the purpose of the exercise
	c.Assert(an.EpCount, Equals, uint32(3)) // NEED >= 2
	c.Assert(an.ClusterSize, Equals, K)

	anID := an.MemberID
	c.Assert(reg.IDCount(), Equals, uint(3)) // regID + anID + clusterID
	clusterID = an.ClusterID

	// DEBUG
	fmt.Printf("regID     %s\n", regID.String())
	fmt.Printf("anID      %s\n", anID.String())
	fmt.Printf("clusterID %s\n", an.ClusterID.String())
	fmt.Printf("  size    %d\n", an.ClusterSize)
	fmt.Printf("  name    %s\n", an.ClusterName)
	// END

	found, err := reg.ContainsID(regID)
	c.Assert(err, IsNil)
	c.Assert(found, Equals, true)

	found, err = reg.ContainsID(anID)
	c.Assert(err, IsNil)
	c.Assert(found, Equals, true)

	found, err = reg.ContainsID(an.ClusterID)
	c.Assert(err, IsNil)
	c.Assert(found, Equals, true)

	return
}

// Create PktComm layers for K members
//
///////////////////////////////////////////////////////////////////////
// XXX epCount IS SET TO 3, WHICH IS WRONG IN GENERAL.
///////////////////////////////////////////////////////////////////////
func (s *XLSuite) createKMemberPktLayers(c *C, rng *xr.PRNG,
	server *xg.RegServer,
	clusterName string, clusterAttrs uint64, clusterID *xi.NodeID,
	K uint32) (pl []*PktLayer, plNames []string) {

	serverName := server.GetName()
	serverID := server.GetNodeID()
	serverEnd := server.GetEndPoint(0)
	serverCK := server.GetCommsPublicKey()
	serverSK := server.GetSigPublicKey()
	c.Assert(serverEnd, NotNil)

	var err error
	pl = make([]*PktLayer, K)
	plNames = make([]string, K)
	namesInUse := make(map[string]bool)
	for i := uint32(0); i < K; i++ {
		var ep *xt.TcpEndPoint
		ep, err = xt.NewTcpEndPoint("127.0.0.1:0")
		c.Assert(err, IsNil)
		e := []xt.EndPointI{ep}
		newName := rng.NextFileName(8)
		_, ok := namesInUse[newName]
		for ok {
			newName = rng.NextFileName(8)
			_, ok = namesInUse[newName]
		}
		namesInUse[newName] = true
		plNames[i] = newName // guaranteed to be LOCALLY unique
		attrs := uint64(rng.Int63())
		pl[i], err = NewPktLayer(plNames[i], "",
			nil, nil, // private RSA keys are generated if nil
			attrs,
			serverName, serverID, serverEnd, serverCK, serverSK,
			clusterName, clusterAttrs, clusterID,
			K, uint32(3), e) //3 is endPoint count
		c.Assert(err, IsNil)
		c.Assert(pl[i], NotNil)
		c.Assert(pl[i].ClusterID, NotNil)
	}
	return
}

// THIS IS INCOMPLETE

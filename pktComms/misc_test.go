package pktComms

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

/////////////////////////////////////////////////////////////////////
// NEW CODE ABOVE THIS LINE, OLD CODE BELOW
/////////////////////////////////////////////////////////////////////

/////////////////////////////////////////////////////////////////////
// These test functions were hacked from xlReg_go and then modified for
// use in testing paxos_go.
/////////////////////////////////////////////////////////////////////

//// Returns a 32-byte slice of random values.  Strength is determined by
//// the strength of rng.
////
//func (s *XLSuite) makeAnID(c *C, rng *xr.PRNG) (id []byte) {
//	id = rng.SomeBytes(SHA3_LEN)
//	return
//}
//
//// Returns a 32-byte NodeID.
//func (s *XLSuite) makeANodeID(c *C, rng *xr.PRNG) (nodeID *xi.NodeID) {
//	id := s.makeAnID(c, rng)
//	nodeID, err := xi.New(id)
//	c.Assert(err, IsNil)
//	c.Assert(nodeID, Not(IsNil))
//	return
//}
//
//// Returns a 2048-bit RSA private key.
//func (s *XLSuite) makeAnRSAKey(c *C) (key *rsa.PrivateKey) {
//	key, err := rsa.GenerateKey(rand.Reader, 2048)
//	c.Assert(err, IsNil)
//	c.Assert(key, Not(IsNil))
//	return key
//}
//
//// Creates a local (127.0.0.1) tcp/ip endPoint and adds it to the node.
//// This code was hacked from xlNode_go/node_test.go.
////
//func (s *XLSuite) makeALocalTCPEndPoint(c *C, node *xn.Node) {
//	addr := fmt.Sprintf("127.0.0.1:0")
//	ep, err := xt.NewTcpEndPoint(addr)
//	c.Assert(err, IsNil)
//	c.Assert(ep, Not(IsNil))
//	ndx, err := node.AddEndPoint(ep)
//	c.Assert(err, IsNil)
//	c.Assert(ndx, Equals, 0) // it's the only one
//}
//
//// Return an initialized and tested host, with a NodeID, commsKey,
//// and sigKey.   This code was hacked from xlNode_go/node_test.go
//// and then simplified a bit.
////
//func (s *XLSuite) makeNodeAndKeys(c *C, rng *xr.PRNG,
//	namesInUse map[string]bool) (n *xn.Node, ckPriv, skPriv *rsa.PrivateKey) {
//
//	var name string
//	for {
//		name = rng.NextFileName(8)
//		for {
//			first := string(name[0])
//			if !strings.Contains(first, "0123456789") &&
//				!strings.Contains(name, "-") {
//				break
//			}
//			name = rng.NextFileName(8)
//		}
//		if _, ok := namesInUse[name]; !ok {
//			// it's not in use
//			namesInUse[name] = true
//			break
//		}
//	}
//	id := s.makeANodeID(c, rng)
//	lfs := "tmp/" + hex.EncodeToString(id.Value())
//	ckPriv = s.makeAnRSAKey(c)
//	skPriv = s.makeAnRSAKey(c)
//
//	n, err2 := xn.New(name, id, lfs, ckPriv, skPriv, nil, nil, nil)
//
//	c.Assert(err2, IsNil)
//	c.Assert(n, Not(IsNil))
//	c.Assert(name, Equals, n.GetName())
//	actualID := n.GetNodeID()
//	c.Assert(true, Equals, id.Equal(actualID))
//	// s.doKeyTests(c, n, rng)
//	c.Assert(0, Equals, (*n).SizePeers())
//	c.Assert(0, Equals, (*n).SizeOverlays())
//	c.Assert(0, Equals, n.SizeConnections())
//	c.Assert(lfs, Equals, n.GetLFS())
//	return n, ckPriv, skPriv
//}
//
//// Create a quasi-random base node and from that a MemberInfo data
//// strplture.
////
//// XXX Using functions must check to ensure members have unique names
////
//func (s *XLSuite) makeAMemberInfo(c *C, rng *xr.PRNG) *xg.MemberInfo {
//	attrs := uint64(rng.Int63())
//	bn, err := xn.NewBaseNode(
//		rng.NextFileName(8),
//		s.makeANodeID(c, rng),
//		&s.makeAnRSAKey(c).PublicKey,
//		&s.makeAnRSAKey(c).PublicKey,
//		nil) // overlays
//	c.Assert(err, IsNil)
//	return &xg.MemberInfo{
//		Attrs:    attrs,
//		BaseNode: *bn,
//	}
//}
//
//// Returns a MemberInfo strplture for a given node.  The BaseNode is
//// cloned.
//func (s *XLSuite) memberInfoForNode(c *C, rng *xr.PRNG, node *xn.Node) *xg.MemberInfo {
//
//	attrs := uint64(rng.Int63())
//	bn, err := xn.NewBaseNode(
//		node.GetName(),
//		node.GetNodeID(),
//		node.GetCommsPublicKey(),
//		node.GetSigPublicKey(),
//		nil) // overlays
//	c.Assert(err, IsNil)
//	return &xg.MemberInfo{
//		Attrs:    attrs,
//		BaseNode: *bn,
//	}
//}
//
//// Make a RegCluster for test purposes.  Cluster member names are guaranteed
//// to be unique but the name of the cluster itself may not be.
////
//func (s *XLSuite) makeACluster(c *C, rng *xr.PRNG, epCount, size uint32) (
//	rc *xg.RegCluster,
//	members []*xn.Node, ckPrivs, skPrivs []*rsa.PrivateKey) {
//
//	var err error
//	c.Assert(xg.MIN_CLUSTER_SIZE <= size && size <= xg.MAX_CLUSTER_SIZE,
//		Equals, true)
//
//	attrs := uint64(rng.Int63())
//	name := rng.NextFileName(8) // no guarantee of uniqueness
//	id := s.makeANodeID(c, rng)
//
//	rc, err = xg.NewRegCluster(name, id, attrs, size, epCount)
//	c.Assert(err, IsNil)
//
//	namesInUse := make(map[string]bool)
//
//	for count := uint32(0); count < size; count++ {
//		member, ckPriv, skPriv := s.makeNodeAndKeys(c, rng, namesInUse)
//		members = append(members, member)
//		ckPrivs = append(ckPrivs, ckPriv)
//		skPrivs = append(skPrivs, skPriv)
//		cm := s.memberInfoForNode(c, rng, member)
//		err = rc.AddMember(cm)
//		c.Assert(err, IsNil)
//
//	}
//	return
//}

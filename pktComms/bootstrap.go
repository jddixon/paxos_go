package pktComms

// paxos_go/pkt_comms/bootstrap.go

import (
	"crypto/rsa"
	"fmt"
	xi "github.com/jddixon/xlNodeID_go"
	xg "github.com/jddixon/xlReg_go"
	xt "github.com/jddixon/xlTransport_go"
	"sync"
)

var _ = fmt.Print

// The bootstrap process enables the caller to join a cluster
// and learn information about the cluster's other members.  Once the
// bootstrapper has learned that information, it is done.

// As implemented so far, this is an ephemeral client, meaning that it
// neither saves nor restores its Node; keys and such are generated for
// each instance.

// For practical use, it is essential that the PktLayer create its
// Node when NewPktLayer() is first called, but then save its
// configuration.  This is conventionally written to LFS/.xlattice/config.
// On subsequent the client reads its configuration file rather than
// regenerating keys, etc.

type Bootstrapper struct {
	DoneCh chan error
	mu sync.RWMutex		// ever used ??? XXX
	xg.MemberNode
}

func NewBootstrapper(
	name, lfs string, ckPriv, skPriv *rsa.PrivateKey, attrs uint64,
	serverName string, serverID *xi.NodeID, serverEnd xt.EndPointI,
	serverCK, serverSK *rsa.PublicKey,
	clusterName string, clusterAttrs uint64, clusterID *xi.NodeID,
	size, epCount uint32, e []xt.EndPointI) (bs *Bootstrapper, err error) {

	if lfs == "" {
		attrs |= xg.ATTR_EPHEMERAL
	}
	mn, err := xg.NewMemberNode(name, lfs, ckPriv, skPriv, attrs,
		serverName, serverID, serverEnd, serverCK, serverSK,
		clusterName, clusterAttrs, clusterID, size, epCount, e)

	if err == nil {
		bs = &Bootstrapper{
			DoneCh:		make(chan error),
			MemberNode: *mn,
		}
	}
	return

}

// Join the cluster and get the members list.
func (bs *Bootstrapper) JoinCluster() {

	mn := &bs.MemberNode
	var (
		err      error
		version1 uint32
	)
	// join the cluster and get the members list
	cnx, version2, err := mn.SessionSetup(version1)
	_ = version2 // not yet used
	if err == nil {
		bs.Cnx = cnx
		err = mn.MemberAndOK()
		if err == nil {
			err = mn.JoinAndReply()
			if err == nil {
				err = mn.GetAndMembers() // XXX PANICS
			}
		}
	}
	// DEBUG
	fmt.Printf("JoinCluster for %s done, returning %v\n", 
		bs.GetName(), err)
	// END
	bs.DoneCh <- err
}

// Start the PktLayer running in separate goroutine, so that this function
// is non-blocking.

func (bs *Bootstrapper) Run() {

	mn := &bs.MemberNode

	go func() {
		var err error

		// DEBUG ------------------------------------------
		var nilMembers []int
		for i := 0; i < len(bs.Members); i++ {
			if bs.Members[i] == nil {
				nilMembers = append(nilMembers, i)
			}
		}
		if len(nilMembers) > 0 {
			fmt.Printf("PktLayer.Run() after Get finds nil members: %v\n",
				nilMembers)
		}
		// END --------------------------------------------
		if err == nil {
			err = mn.ByeAndAck()
		}

		// END OF RUN ===============================================
		if bs.Cnx != nil {
			bs.Cnx.Close()
		}
		bs.DoneCh <- err
	}()
}

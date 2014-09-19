package pktComms

// paxos_go/pkt_comms/pktLayer.go

import (
	"crypto/rsa"
	"fmt"
	xi "github.com/jddixon/xlNodeID_go"
	xg "github.com/jddixon/xlReg_go"
	xt "github.com/jddixon/xlTransport_go"
	"sync"
)

var _ = fmt.Print

// The PktLayer is created to enable the caller to join a cluster
// and learn information about the cluster's other members.  Once the
// client has learned that information, it is done.

// As implemented so far, this is an ephemeral client, meaning that it
// neither saves nor restores its Node; keys and such are generated for
// each instance.

// For practical use, it is essential that the PktLayer create its
// Node when NewPktLayer() is first called, but then save its
// configuration.  This is conventionally written to LFS/.xlattice/config.
// On subsequent the client reads its configuration file rather than
// regenerating keys, etc.

type PktLayer struct {
	Cnx *xt.TcpConnection
	mu  sync.RWMutex
	xg.MemberNode
}

func NewPktLayer(
	name, lfs string, ckPriv, skPriv *rsa.PrivateKey, attrs uint64,
	serverName string, serverID *xi.NodeID, serverEnd xt.EndPointI,
	serverCK, serverSK *rsa.PublicKey,
	clusterName string, clusterAttrs uint64, clusterID *xi.NodeID,
	size, epCount uint32, e []xt.EndPointI) (pl *PktLayer, err error) {

	if lfs == "" {
		attrs |= xg.ATTR_EPHEMERAL
	}
	mn, err := xg.NewMemberNode(name, lfs, ckPriv, skPriv, attrs,
		serverName, serverID, serverEnd, serverCK, serverSK,
		clusterName, clusterAttrs, clusterID, size, epCount, e)

	if err == nil {
		mn.DoneCh = make(chan bool)
		pl = &PktLayer{
			MemberNode: *mn,
		}
	}
	return

}

// Join the cluster and get the members list.
func (pl *PktLayer) JoinCluster() {

	mn := &pl.MemberNode
	var (
		err      error
		version1 uint32
	)
	// join the cluster and get the members list
	cnx, version2, err := mn.SessionSetup(version1)
	_ = version2 // not yet used
	if err == nil {
		pl.Cnx = cnx
		err = mn.MemberAndOK()
		if err == nil {
			err = mn.JoinAndReply()
			if err == nil {
				err = mn.GetAndMembers() // XXX PANICS
			}
		}
	}
	if err == nil {
		pl.DoneCh <- true
	} else {
		pl.Err = err
		pl.DoneCh <- false
	}
}

// Start the PktLayer running in separate goroutine, so that this function
// is non-blocking.

func (pl *PktLayer) Run() {

	mn := &pl.MemberNode
	// pl.JoinCluster() // XXX SHOULD NOT BE HERE

	go func() {
		var err error

		// DEBUG ------------------------------------------
		var nilMembers []int
		for i := 0; i < len(pl.Members); i++ {
			if pl.Members[i] == nil {
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
		if pl.Cnx != nil {
			pl.Cnx.Close()
		}
		if err != nil {
			mn.Err = err
			mn.DoneCh <- false
		} else {
			mn.DoneCh <- true
		}
	}()
	return
}

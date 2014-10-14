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


type PktLayer struct {
	Cnx *xt.TcpConnection
	mu  sync.RWMutex
	xg.MemberNode
}

// XXX This is WRONG.  Main argument is *xg.MemberNode, which we get
// by deserializing what's in LFS/.xlattice

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
		pl = &PktLayer{
			MemberNode: *mn,
		}
	}
	return

}

// Start the PktLayer running in separate goroutine, so that this function
// is non-blocking.

func (pl *PktLayer) Run() {

	mn := &pl.MemberNode

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
		mn.DoneCh <- err
	}()
}

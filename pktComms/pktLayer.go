package pktComms

// paxos_go/pkt_comms/pktLayer.go

import (
	//"crypto/rsa"
	"fmt"
	//xi "github.com/jddixon/xlNodeID_go"
	xg "github.com/jddixon/xlReg_go"
	xt "github.com/jddixon/xlTransport_go"
	"sync"
)

var _ = fmt.Print


type PktLayer struct {
	StopCh		chan bool
	StoppedCh	chan error
	Cnx *xt.TcpConnection		// value?
	mu  sync.RWMutex
	PktCommsNode
}

func NewPktLayer(o *PktLayerOptions) (pl *PktLayer, err error) {

	if o.LFS == "" {
		o.Attrs |= xg.ATTR_EPHEMERAL
	}
	// XXX HACK TO MAKE THINGS COMPILE
	mn, err := xg.NewMemberNode(
		o.Name, o.LFS, o.CKPriv, o.SKPriv, o.Attrs, 
		o.ServerName, o.ServerID, o.ServerEnd, o.ServerCK, o.ServerSK,
		o.ClusterName, o.ClusterAttrs, o.ClusterID, o.Size,
		o.EPCount, o.EndPoints)

	if err == nil {
		pcn := &PktCommsNode{
			MemberNode: *mn,
		}
		pl = &PktLayer{
			StopCh		: make(chan bool),
			StoppedCh	: make(chan error),
			PktCommsNode: *pcn,
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

package pktComms

// paxos_go/pkt_comms/pktLayer.go

import (
	"fmt"
	xcl "github.com/jddixon/xlCluster_go"
	//xi "github.com/jddixon/xlNodeID_go"
	//xn "github.com/jddixon/xlNode_go"
	xg "github.com/jddixon/xlReg_go"
	//xt "github.com/jddixon/xlTransport_go"
	"sync"
)

var _ = fmt.Print

type PktLayer struct {
	DoneCh chan error
	// VALUE?
	StopCh    chan bool
	StoppedCh chan error
	// END VALUE?
	Running bool
	Mu      sync.RWMutex
	PktCommsNode
}

func NewPktLayer(cm *xcl.ClusterMember) (pl *PktLayer, err error) {

	if cm == nil {
		err = NilClusterMember
	} else {
		lfs := cm.GetLFS()
		if lfs == "" {
			cm.Attrs |= xg.ATTR_EPHEMERAL
		}
		pcn := &PktCommsNode{
			ClusterMember: *cm,
		}
		pl = &PktLayer{
			DoneCh: make(chan error),
			// VALUE?
			StopCh:    make(chan bool),
			StoppedCh: make(chan error),
			// END VALUE?
			PktCommsNode: *pcn,
		}
	}
	return

}

// Start the PktLayer running in separate goroutine, so that this function
// is non-blocking.

func (pl *PktLayer) Start() {

	cl := &pl.ClusterMember
	err := cl.Run()

	if err == nil {
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
				fmt.Printf("PktLayer.Start() after Get finds nil members: %v\n",
					nilMembers)
			}
			// END --------------------------------------------
			if err == nil {
				// err = mn.ByeAndAck()
			}

			// END OF RUN ===============================================
			pl.DoneCh <- err
		}()
	} else {
		pl.DoneCh <- err
	}
}

package pktComms

// paxos_go/pkt_comms/pktCommsNode.go

import (
	xcl "github.com/jddixon/xlCluster_go"
	//xg "github.com/jddixon/xlReg_go"
)

type PktCommsNode struct {
	xcl.ClusterMember
}

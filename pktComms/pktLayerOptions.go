package pktComms

// paxos_go/pkt_comms/pktLayerOptions.go

import (
	"crypto/rsa"
	"fmt"
	xi "github.com/jddixon/xlNodeID_go"
	//xg "github.com/jddixon/xlReg_go"
	xt "github.com/jddixon/xlTransport_go"
)

var _ = fmt.Print

// Use this and then null it to encourage GC.
type PktLayerOptions struct {
	Rebooting bool
	Name string
	LFS string 
	CKPriv *rsa.PrivateKey
	SKPriv *rsa.PrivateKey
	Attrs uint64
	ServerName string
	ServerID *xi.NodeID
	ServerEnd xt.EndPointI
	ServerCK *rsa.PublicKey
	ServerSK *rsa.PublicKey
	ClusterName string
	ClusterAttrs uint64
	ClusterID *xi.NodeID
	Size uint32
	EPCount uint32
	EndPoints []xt.EndPointI
}

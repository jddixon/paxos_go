package pktComms

import (
	"bytes"
	"code.google.com/p/go.crypto/sha3"
	"code.google.com/p/goprotobuf/proto"
	"encoding/binary"
	"fmt"
	xr "github.com/jddixon/rnglib_go"
	. "gopkg.in/check.v1"
)

var _ = proto.Marshal
var _ = fmt.Println

func (d *XLSuite) TestPaxosPkt(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("TEST_PAXOS_PKT")
	}

	rng := xr.MakeSimpleRNG()

	mySeqN := uint64(rng.Int63())
	for mySeqN == 0 { // may not be zero
		mySeqN = uint64(rng.Int63())
	}

	id := rng.SomeBytes(SHA3_LEN)

	seqBuf := new(bytes.Buffer)
	binary.Write(seqBuf, binary.LittleEndian, mySeqN)

	msgLen := uint(64 + rng.Intn(64))
	msg	:= rng.SomeBytes(msgLen) // fill with rubbish
	salt := rng.SomeBytes(8)		// still more rubbish

	digest := sha3.NewKeccak256()
	digest.Write(id)
	digest.Write(seqBuf.Bytes())
	digest.Write(msg)
	digest.Write([]byte(salt))

	hash := digest.Sum(nil)

	var pkt = AppMsg{MsgN: &mySeqN,
		Contents: msg,
		Salt:     salt,
		Hash:     hash}

	// In each of these cases, the test proves that the field
	// was present; otherwise the 'empty' value (zero, nil, etc)
	// would have been returned.
	seqNOut := pkt.GetMsgN()
	c.Assert(seqNOut, Equals, mySeqN)

	msgOut := pkt.GetContents()
	// gocheck can't compare byte arrays
	// c.Assert( msgOut, Equals, msg)
	d.compareByteSlices(c, msgOut, msg)

	saltOut := pkt.GetSalt()
	d.compareByteSlices(c, saltOut, salt)

	hashOut := pkt.GetHash()
	d.compareByteSlices(c, hashOut, hash)
}

func (d *XLSuite) compareByteSlices(c *C, a []byte, b []byte) {
	c.Assert(len(a), Equals, len(b))
	for i := 0; i < len(b); i++ {
		c.Assert(a[i], Equals, b[i])
	}
}

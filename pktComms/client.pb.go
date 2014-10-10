// Code generated by protoc-gen-go.
// source: client.proto
// DO NOT EDIT!

package pktComms

import proto "code.google.com/p/goprotobuf/proto"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = math.Inf

type Request struct {
	ClientID         *uint64 `protobuf:"varint,1,opt" json:"ClientID,omitempty"`
	RequestID        []byte  `protobuf:"bytes,2,opt" json:"RequestID,omitempty"`
	Cmd              []byte  `protobuf:"bytes,3,opt" json:"Cmd,omitempty"`
	Salt             []byte  `protobuf:"bytes,12,opt" json:"Salt,omitempty"`
	DigSig           []byte  `protobuf:"bytes,13,opt" json:"DigSig,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *Request) Reset()         { *m = Request{} }
func (m *Request) String() string { return proto.CompactTextString(m) }
func (*Request) ProtoMessage()    {}

func (m *Request) GetClientID() uint64 {
	if m != nil && m.ClientID != nil {
		return *m.ClientID
	}
	return 0
}

func (m *Request) GetRequestID() []byte {
	if m != nil {
		return m.RequestID
	}
	return nil
}

func (m *Request) GetCmd() []byte {
	if m != nil {
		return m.Cmd
	}
	return nil
}

func (m *Request) GetSalt() []byte {
	if m != nil {
		return m.Salt
	}
	return nil
}

func (m *Request) GetDigSig() []byte {
	if m != nil {
		return m.DigSig
	}
	return nil
}

type Response struct {
	ClientID         *uint64 `protobuf:"varint,1,opt" json:"ClientID,omitempty"`
	RequestID        []byte  `protobuf:"bytes,2,opt" json:"RequestID,omitempty"`
	Resp             []byte  `protobuf:"bytes,3,opt" json:"Resp,omitempty"`
	Salt             []byte  `protobuf:"bytes,12,opt" json:"Salt,omitempty"`
	DigSig           []byte  `protobuf:"bytes,13,opt" json:"DigSig,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *Response) Reset()         { *m = Response{} }
func (m *Response) String() string { return proto.CompactTextString(m) }
func (*Response) ProtoMessage()    {}

func (m *Response) GetClientID() uint64 {
	if m != nil && m.ClientID != nil {
		return *m.ClientID
	}
	return 0
}

func (m *Response) GetRequestID() []byte {
	if m != nil {
		return m.RequestID
	}
	return nil
}

func (m *Response) GetResp() []byte {
	if m != nil {
		return m.Resp
	}
	return nil
}

func (m *Response) GetSalt() []byte {
	if m != nil {
		return m.Salt
	}
	return nil
}

func (m *Response) GetDigSig() []byte {
	if m != nil {
		return m.DigSig
	}
	return nil
}

type Redirect struct {
	ClientID         *uint64 `protobuf:"varint,1,opt" json:"ClientID,omitempty"`
	RequestID        []byte  `protobuf:"bytes,2,opt" json:"RequestID,omitempty"`
	Resp             *string `protobuf:"bytes,3,opt" json:"Resp,omitempty"`
	Salt             []byte  `protobuf:"bytes,12,opt" json:"Salt,omitempty"`
	DigSig           []byte  `protobuf:"bytes,13,opt" json:"DigSig,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *Redirect) Reset()         { *m = Redirect{} }
func (m *Redirect) String() string { return proto.CompactTextString(m) }
func (*Redirect) ProtoMessage()    {}

func (m *Redirect) GetClientID() uint64 {
	if m != nil && m.ClientID != nil {
		return *m.ClientID
	}
	return 0
}

func (m *Redirect) GetRequestID() []byte {
	if m != nil {
		return m.RequestID
	}
	return nil
}

func (m *Redirect) GetResp() string {
	if m != nil && m.Resp != nil {
		return *m.Resp
	}
	return ""
}

func (m *Redirect) GetSalt() []byte {
	if m != nil {
		return m.Salt
	}
	return nil
}

func (m *Redirect) GetDigSig() []byte {
	if m != nil {
		return m.DigSig
	}
	return nil
}

func init() {
}
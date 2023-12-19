// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: system.proto

package system

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type Session struct {
	ID                   string                   `protobuf:"bytes,1,opt,name=ID,proto3" json:"ID,omitempty"`
	Topics               map[string]*SessionTopic `protobuf:"bytes,2,rep,name=Topics,proto3" json:"Topics,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	LastActiveTime       int64                    `protobuf:"varint,3,opt,name=LastActiveTime,proto3" json:"LastActiveTime,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                 `json:"-"`
	XXX_unrecognized     []byte                   `json:"-"`
	XXX_sizecache        int32                    `json:"-"`
}

func (m *Session) Reset()         { *m = Session{} }
func (m *Session) String() string { return proto.CompactTextString(m) }
func (*Session) ProtoMessage()    {}
func (*Session) Descriptor() ([]byte, []int) {
	return fileDescriptor_86a7260ebdc12f47, []int{0}
}
func (m *Session) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Session) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Session.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Session) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Session.Merge(m, src)
}
func (m *Session) XXX_Size() int {
	return m.Size()
}
func (m *Session) XXX_DiscardUnknown() {
	xxx_messageInfo_Session.DiscardUnknown(m)
}

var xxx_messageInfo_Session proto.InternalMessageInfo

func (m *Session) GetID() string {
	if m != nil {
		return m.ID
	}
	return ""
}

func (m *Session) GetTopics() map[string]*SessionTopic {
	if m != nil {
		return m.Topics
	}
	return nil
}

func (m *Session) GetLastActiveTime() int64 {
	if m != nil {
		return m.LastActiveTime
	}
	return 0
}

type SessionTopic struct {
	Topic                string   `protobuf:"bytes,1,opt,name=topic,proto3" json:"topic,omitempty"`
	QoS                  int32    `protobuf:"varint,2,opt,name=QoS,proto3" json:"QoS,omitempty"`
	UnAckMessageID       []string `protobuf:"bytes,3,rep,name=UnAckMessageID,proto3" json:"UnAckMessageID,omitempty"`
	UnReceivePacketID    []uint32 `protobuf:"varint,4,rep,packed,name=UnReceivePacketID,proto3" json:"UnReceivePacketID,omitempty"`
	UnCompletePacketID   []uint32 `protobuf:"varint,5,rep,packed,name=UnCompletePacketID,proto3" json:"UnCompletePacketID,omitempty"`
	LastAckMessageID     string   `protobuf:"bytes,6,opt,name=LastAckMessageID,proto3" json:"LastAckMessageID,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SessionTopic) Reset()         { *m = SessionTopic{} }
func (m *SessionTopic) String() string { return proto.CompactTextString(m) }
func (*SessionTopic) ProtoMessage()    {}
func (*SessionTopic) Descriptor() ([]byte, []int) {
	return fileDescriptor_86a7260ebdc12f47, []int{1}
}
func (m *SessionTopic) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *SessionTopic) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_SessionTopic.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *SessionTopic) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SessionTopic.Merge(m, src)
}
func (m *SessionTopic) XXX_Size() int {
	return m.Size()
}
func (m *SessionTopic) XXX_DiscardUnknown() {
	xxx_messageInfo_SessionTopic.DiscardUnknown(m)
}

var xxx_messageInfo_SessionTopic proto.InternalMessageInfo

func (m *SessionTopic) GetTopic() string {
	if m != nil {
		return m.Topic
	}
	return ""
}

func (m *SessionTopic) GetQoS() int32 {
	if m != nil {
		return m.QoS
	}
	return 0
}

func (m *SessionTopic) GetUnAckMessageID() []string {
	if m != nil {
		return m.UnAckMessageID
	}
	return nil
}

func (m *SessionTopic) GetUnReceivePacketID() []uint32 {
	if m != nil {
		return m.UnReceivePacketID
	}
	return nil
}

func (m *SessionTopic) GetUnCompletePacketID() []uint32 {
	if m != nil {
		return m.UnCompletePacketID
	}
	return nil
}

func (m *SessionTopic) GetLastAckMessageID() string {
	if m != nil {
		return m.LastAckMessageID
	}
	return ""
}

type SystemState struct {
	Session              *Session `protobuf:"bytes,1,opt,name=Session,proto3" json:"Session,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SystemState) Reset()         { *m = SystemState{} }
func (m *SystemState) String() string { return proto.CompactTextString(m) }
func (*SystemState) ProtoMessage()    {}
func (*SystemState) Descriptor() ([]byte, []int) {
	return fileDescriptor_86a7260ebdc12f47, []int{2}
}
func (m *SystemState) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *SystemState) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_SystemState.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *SystemState) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SystemState.Merge(m, src)
}
func (m *SystemState) XXX_Size() int {
	return m.Size()
}
func (m *SystemState) XXX_DiscardUnknown() {
	xxx_messageInfo_SystemState.DiscardUnknown(m)
}

var xxx_messageInfo_SystemState proto.InternalMessageInfo

func (m *SystemState) GetSession() *Session {
	if m != nil {
		return m.Session
	}
	return nil
}

func init() {
	proto.RegisterType((*Session)(nil), "system.Session")
	proto.RegisterMapType((map[string]*SessionTopic)(nil), "system.Session.TopicsEntry")
	proto.RegisterType((*SessionTopic)(nil), "system.SessionTopic")
	proto.RegisterType((*SystemState)(nil), "system.SystemState")
}

func init() { proto.RegisterFile("system.proto", fileDescriptor_86a7260ebdc12f47) }

var fileDescriptor_86a7260ebdc12f47 = []byte{
	// 350 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x52, 0xc1, 0x4a, 0xf3, 0x40,
	0x10, 0xfe, 0x37, 0xf9, 0x93, 0x9f, 0x4e, 0xfa, 0xb7, 0x75, 0xe9, 0x21, 0x28, 0x84, 0xd8, 0x83,
	0xc4, 0x22, 0x11, 0xda, 0x4b, 0xd1, 0x53, 0x35, 0x1e, 0x02, 0x8a, 0xba, 0x69, 0x2f, 0xde, 0x62,
	0x18, 0x24, 0xb4, 0x4d, 0x4a, 0x77, 0x2d, 0xf4, 0x4d, 0x7c, 0x23, 0x3d, 0xfa, 0x08, 0x52, 0x2f,
	0x3e, 0x86, 0x64, 0x37, 0x6a, 0x49, 0x3d, 0xed, 0xcc, 0x7c, 0xdf, 0xcc, 0x7c, 0xfb, 0xed, 0x42,
	0x9d, 0xaf, 0xb8, 0xc0, 0x99, 0x3f, 0x5f, 0xe4, 0x22, 0xa7, 0xa6, 0xca, 0x3a, 0xcf, 0x04, 0xfe,
	0x45, 0xc8, 0x79, 0x9a, 0x67, 0xb4, 0x01, 0x5a, 0x18, 0xd8, 0xc4, 0x25, 0x5e, 0x8d, 0x69, 0x61,
	0x40, 0xfb, 0x60, 0x8e, 0xf2, 0x79, 0x9a, 0x70, 0x5b, 0x73, 0x75, 0xcf, 0xea, 0xed, 0xf9, 0xe5,
	0x88, 0xb2, 0xc1, 0x57, 0xe8, 0x45, 0x26, 0x16, 0x2b, 0x56, 0x52, 0xe9, 0x01, 0x34, 0x2e, 0x63,
	0x2e, 0x86, 0x89, 0x48, 0x97, 0x38, 0x4a, 0x67, 0x68, 0xeb, 0x2e, 0xf1, 0x74, 0x56, 0xa9, 0xee,
	0x5e, 0x83, 0xb5, 0xd1, 0x4e, 0x5b, 0xa0, 0x4f, 0x70, 0x55, 0x2e, 0x2f, 0x42, 0xda, 0x05, 0x63,
	0x19, 0x4f, 0x1f, 0xd1, 0xd6, 0x5c, 0xe2, 0x59, 0xbd, 0x76, 0x65, 0xb9, 0x6c, 0x66, 0x8a, 0x72,
	0xa2, 0x0d, 0x48, 0xe7, 0x83, 0x40, 0x7d, 0x13, 0xa3, 0x6d, 0x30, 0x64, 0x50, 0x0e, 0x55, 0x49,
	0xb1, 0xe8, 0x36, 0x8f, 0xe4, 0x50, 0x83, 0x15, 0x61, 0xa1, 0x78, 0x9c, 0x0d, 0x93, 0xc9, 0x15,
	0x72, 0x1e, 0x3f, 0x60, 0x18, 0xd8, 0xba, 0xab, 0x7b, 0x35, 0x56, 0xa9, 0xd2, 0x23, 0xd8, 0x19,
	0x67, 0x0c, 0x13, 0x4c, 0x97, 0x78, 0x13, 0x27, 0x13, 0x14, 0x61, 0x60, 0xff, 0x75, 0x75, 0xef,
	0x3f, 0xdb, 0x06, 0xa8, 0x0f, 0x74, 0x9c, 0x9d, 0xe7, 0xb3, 0xf9, 0x14, 0xc5, 0x0f, 0xdd, 0x90,
	0xf4, 0x5f, 0x10, 0xda, 0x85, 0x96, 0x72, 0x68, 0x43, 0x87, 0x29, 0x85, 0x6f, 0xd5, 0x3b, 0x03,
	0xb0, 0x22, 0x69, 0x46, 0x24, 0x62, 0x81, 0xf4, 0xf0, 0xfb, 0x09, 0xe5, 0x55, 0xad, 0x5e, 0xb3,
	0xe2, 0x15, 0xfb, 0xc2, 0xcf, 0xf6, 0x5f, 0xd6, 0x0e, 0x79, 0x5d, 0x3b, 0xe4, 0x6d, 0xed, 0x90,
	0xa7, 0x77, 0xe7, 0xcf, 0x5d, 0xd3, 0x3f, 0x56, 0xe4, 0x53, 0x75, 0xdc, 0x9b, 0xf2, 0x83, 0xf4,
	0x3f, 0x03, 0x00, 0x00, 0xff, 0xff, 0xa5, 0x40, 0xc5, 0x63, 0x30, 0x02, 0x00, 0x00,
}

func (m *Session) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Session) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Session) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if m.LastActiveTime != 0 {
		i = encodeVarintSystem(dAtA, i, uint64(m.LastActiveTime))
		i--
		dAtA[i] = 0x18
	}
	if len(m.Topics) > 0 {
		for k := range m.Topics {
			v := m.Topics[k]
			baseI := i
			if v != nil {
				{
					size, err := v.MarshalToSizedBuffer(dAtA[:i])
					if err != nil {
						return 0, err
					}
					i -= size
					i = encodeVarintSystem(dAtA, i, uint64(size))
				}
				i--
				dAtA[i] = 0x12
			}
			i -= len(k)
			copy(dAtA[i:], k)
			i = encodeVarintSystem(dAtA, i, uint64(len(k)))
			i--
			dAtA[i] = 0xa
			i = encodeVarintSystem(dAtA, i, uint64(baseI-i))
			i--
			dAtA[i] = 0x12
		}
	}
	if len(m.ID) > 0 {
		i -= len(m.ID)
		copy(dAtA[i:], m.ID)
		i = encodeVarintSystem(dAtA, i, uint64(len(m.ID)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *SessionTopic) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *SessionTopic) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *SessionTopic) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.LastAckMessageID) > 0 {
		i -= len(m.LastAckMessageID)
		copy(dAtA[i:], m.LastAckMessageID)
		i = encodeVarintSystem(dAtA, i, uint64(len(m.LastAckMessageID)))
		i--
		dAtA[i] = 0x32
	}
	if len(m.UnCompletePacketID) > 0 {
		dAtA3 := make([]byte, len(m.UnCompletePacketID)*10)
		var j2 int
		for _, num := range m.UnCompletePacketID {
			for num >= 1<<7 {
				dAtA3[j2] = uint8(uint64(num)&0x7f | 0x80)
				num >>= 7
				j2++
			}
			dAtA3[j2] = uint8(num)
			j2++
		}
		i -= j2
		copy(dAtA[i:], dAtA3[:j2])
		i = encodeVarintSystem(dAtA, i, uint64(j2))
		i--
		dAtA[i] = 0x2a
	}
	if len(m.UnReceivePacketID) > 0 {
		dAtA5 := make([]byte, len(m.UnReceivePacketID)*10)
		var j4 int
		for _, num := range m.UnReceivePacketID {
			for num >= 1<<7 {
				dAtA5[j4] = uint8(uint64(num)&0x7f | 0x80)
				num >>= 7
				j4++
			}
			dAtA5[j4] = uint8(num)
			j4++
		}
		i -= j4
		copy(dAtA[i:], dAtA5[:j4])
		i = encodeVarintSystem(dAtA, i, uint64(j4))
		i--
		dAtA[i] = 0x22
	}
	if len(m.UnAckMessageID) > 0 {
		for iNdEx := len(m.UnAckMessageID) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.UnAckMessageID[iNdEx])
			copy(dAtA[i:], m.UnAckMessageID[iNdEx])
			i = encodeVarintSystem(dAtA, i, uint64(len(m.UnAckMessageID[iNdEx])))
			i--
			dAtA[i] = 0x1a
		}
	}
	if m.QoS != 0 {
		i = encodeVarintSystem(dAtA, i, uint64(m.QoS))
		i--
		dAtA[i] = 0x10
	}
	if len(m.Topic) > 0 {
		i -= len(m.Topic)
		copy(dAtA[i:], m.Topic)
		i = encodeVarintSystem(dAtA, i, uint64(len(m.Topic)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *SystemState) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *SystemState) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *SystemState) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if m.Session != nil {
		{
			size, err := m.Session.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintSystem(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintSystem(dAtA []byte, offset int, v uint64) int {
	offset -= sovSystem(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Session) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.ID)
	if l > 0 {
		n += 1 + l + sovSystem(uint64(l))
	}
	if len(m.Topics) > 0 {
		for k, v := range m.Topics {
			_ = k
			_ = v
			l = 0
			if v != nil {
				l = v.Size()
				l += 1 + sovSystem(uint64(l))
			}
			mapEntrySize := 1 + len(k) + sovSystem(uint64(len(k))) + l
			n += mapEntrySize + 1 + sovSystem(uint64(mapEntrySize))
		}
	}
	if m.LastActiveTime != 0 {
		n += 1 + sovSystem(uint64(m.LastActiveTime))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *SessionTopic) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Topic)
	if l > 0 {
		n += 1 + l + sovSystem(uint64(l))
	}
	if m.QoS != 0 {
		n += 1 + sovSystem(uint64(m.QoS))
	}
	if len(m.UnAckMessageID) > 0 {
		for _, s := range m.UnAckMessageID {
			l = len(s)
			n += 1 + l + sovSystem(uint64(l))
		}
	}
	if len(m.UnReceivePacketID) > 0 {
		l = 0
		for _, e := range m.UnReceivePacketID {
			l += sovSystem(uint64(e))
		}
		n += 1 + sovSystem(uint64(l)) + l
	}
	if len(m.UnCompletePacketID) > 0 {
		l = 0
		for _, e := range m.UnCompletePacketID {
			l += sovSystem(uint64(e))
		}
		n += 1 + sovSystem(uint64(l)) + l
	}
	l = len(m.LastAckMessageID)
	if l > 0 {
		n += 1 + l + sovSystem(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *SystemState) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Session != nil {
		l = m.Session.Size()
		n += 1 + l + sovSystem(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func sovSystem(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozSystem(x uint64) (n int) {
	return sovSystem(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Session) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowSystem
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Session: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Session: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ID", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSystem
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthSystem
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthSystem
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ID = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Topics", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSystem
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthSystem
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthSystem
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Topics == nil {
				m.Topics = make(map[string]*SessionTopic)
			}
			var mapkey string
			var mapvalue *SessionTopic
			for iNdEx < postIndex {
				entryPreIndex := iNdEx
				var wire uint64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowSystem
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					wire |= uint64(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				fieldNum := int32(wire >> 3)
				if fieldNum == 1 {
					var stringLenmapkey uint64
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowSystem
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						stringLenmapkey |= uint64(b&0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					intStringLenmapkey := int(stringLenmapkey)
					if intStringLenmapkey < 0 {
						return ErrInvalidLengthSystem
					}
					postStringIndexmapkey := iNdEx + intStringLenmapkey
					if postStringIndexmapkey < 0 {
						return ErrInvalidLengthSystem
					}
					if postStringIndexmapkey > l {
						return io.ErrUnexpectedEOF
					}
					mapkey = string(dAtA[iNdEx:postStringIndexmapkey])
					iNdEx = postStringIndexmapkey
				} else if fieldNum == 2 {
					var mapmsglen int
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowSystem
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						mapmsglen |= int(b&0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					if mapmsglen < 0 {
						return ErrInvalidLengthSystem
					}
					postmsgIndex := iNdEx + mapmsglen
					if postmsgIndex < 0 {
						return ErrInvalidLengthSystem
					}
					if postmsgIndex > l {
						return io.ErrUnexpectedEOF
					}
					mapvalue = &SessionTopic{}
					if err := mapvalue.Unmarshal(dAtA[iNdEx:postmsgIndex]); err != nil {
						return err
					}
					iNdEx = postmsgIndex
				} else {
					iNdEx = entryPreIndex
					skippy, err := skipSystem(dAtA[iNdEx:])
					if err != nil {
						return err
					}
					if (skippy < 0) || (iNdEx+skippy) < 0 {
						return ErrInvalidLengthSystem
					}
					if (iNdEx + skippy) > postIndex {
						return io.ErrUnexpectedEOF
					}
					iNdEx += skippy
				}
			}
			m.Topics[mapkey] = mapvalue
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field LastActiveTime", wireType)
			}
			m.LastActiveTime = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSystem
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.LastActiveTime |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipSystem(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthSystem
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *SessionTopic) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowSystem
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: SessionTopic: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: SessionTopic: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field topic", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSystem
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthSystem
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthSystem
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Topic = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field QoS", wireType)
			}
			m.QoS = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSystem
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.QoS |= int32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field UnAckMessageID", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSystem
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthSystem
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthSystem
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.UnAckMessageID = append(m.UnAckMessageID, string(dAtA[iNdEx:postIndex]))
			iNdEx = postIndex
		case 4:
			if wireType == 0 {
				var v uint32
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowSystem
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					v |= uint32(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				m.UnReceivePacketID = append(m.UnReceivePacketID, v)
			} else if wireType == 2 {
				var packedLen int
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowSystem
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					packedLen |= int(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				if packedLen < 0 {
					return ErrInvalidLengthSystem
				}
				postIndex := iNdEx + packedLen
				if postIndex < 0 {
					return ErrInvalidLengthSystem
				}
				if postIndex > l {
					return io.ErrUnexpectedEOF
				}
				var elementCount int
				var count int
				for _, integer := range dAtA[iNdEx:postIndex] {
					if integer < 128 {
						count++
					}
				}
				elementCount = count
				if elementCount != 0 && len(m.UnReceivePacketID) == 0 {
					m.UnReceivePacketID = make([]uint32, 0, elementCount)
				}
				for iNdEx < postIndex {
					var v uint32
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowSystem
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						v |= uint32(b&0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					m.UnReceivePacketID = append(m.UnReceivePacketID, v)
				}
			} else {
				return fmt.Errorf("proto: wrong wireType = %d for field UnReceivePacketID", wireType)
			}
		case 5:
			if wireType == 0 {
				var v uint32
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowSystem
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					v |= uint32(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				m.UnCompletePacketID = append(m.UnCompletePacketID, v)
			} else if wireType == 2 {
				var packedLen int
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowSystem
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					packedLen |= int(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				if packedLen < 0 {
					return ErrInvalidLengthSystem
				}
				postIndex := iNdEx + packedLen
				if postIndex < 0 {
					return ErrInvalidLengthSystem
				}
				if postIndex > l {
					return io.ErrUnexpectedEOF
				}
				var elementCount int
				var count int
				for _, integer := range dAtA[iNdEx:postIndex] {
					if integer < 128 {
						count++
					}
				}
				elementCount = count
				if elementCount != 0 && len(m.UnCompletePacketID) == 0 {
					m.UnCompletePacketID = make([]uint32, 0, elementCount)
				}
				for iNdEx < postIndex {
					var v uint32
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowSystem
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						v |= uint32(b&0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					m.UnCompletePacketID = append(m.UnCompletePacketID, v)
				}
			} else {
				return fmt.Errorf("proto: wrong wireType = %d for field UnCompletePacketID", wireType)
			}
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field LastAckMessageID", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSystem
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthSystem
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthSystem
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.LastAckMessageID = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipSystem(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthSystem
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *SystemState) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowSystem
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: SystemState: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: SystemState: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Session", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSystem
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthSystem
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthSystem
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Session == nil {
				m.Session = &Session{}
			}
			if err := m.Session.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipSystem(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthSystem
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipSystem(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowSystem
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowSystem
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowSystem
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthSystem
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupSystem
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthSystem
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthSystem        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowSystem          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupSystem = fmt.Errorf("proto: unexpected end of group")
)

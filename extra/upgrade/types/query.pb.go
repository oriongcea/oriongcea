// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: ocea/upgrade/v1beta1/query.proto

package types

import (
	context "context"
	fmt "fmt"
	grpc1 "github.com/gogo/protobuf/grpc"
	proto "github.com/gogo/protobuf/proto"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// QueryCurrentPlanRequest is the request type for the Query/CurrentPlan RPC
// method.
type QueryCurrentPlanRequest struct {
}

func (m *QueryCurrentPlanRequest) Reset()         { *m = QueryCurrentPlanRequest{} }
func (m *QueryCurrentPlanRequest) String() string { return proto.CompactTextString(m) }
func (*QueryCurrentPlanRequest) ProtoMessage()    {}
func (*QueryCurrentPlanRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_4a334d07ad8374f0, []int{0}
}
func (m *QueryCurrentPlanRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryCurrentPlanRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryCurrentPlanRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryCurrentPlanRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryCurrentPlanRequest.Merge(m, src)
}
func (m *QueryCurrentPlanRequest) XXX_Size() int {
	return m.Size()
}
func (m *QueryCurrentPlanRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryCurrentPlanRequest.DiscardUnknown(m)
}

var xxx_messageInfo_QueryCurrentPlanRequest proto.InternalMessageInfo

// QueryCurrentPlanResponse is the response type for the Query/CurrentPlan RPC
// method.
type QueryCurrentPlanResponse struct {
	// plan is the current upgrade plan.
	Plan *Plan `protobuf:"bytes,1,opt,name=plan,proto3" json:"plan,omitempty"`
}

func (m *QueryCurrentPlanResponse) Reset()         { *m = QueryCurrentPlanResponse{} }
func (m *QueryCurrentPlanResponse) String() string { return proto.CompactTextString(m) }
func (*QueryCurrentPlanResponse) ProtoMessage()    {}
func (*QueryCurrentPlanResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_4a334d07ad8374f0, []int{1}
}
func (m *QueryCurrentPlanResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryCurrentPlanResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryCurrentPlanResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryCurrentPlanResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryCurrentPlanResponse.Merge(m, src)
}
func (m *QueryCurrentPlanResponse) XXX_Size() int {
	return m.Size()
}
func (m *QueryCurrentPlanResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryCurrentPlanResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryCurrentPlanResponse proto.InternalMessageInfo

func (m *QueryCurrentPlanResponse) GetPlan() *Plan {
	if m != nil {
		return m.Plan
	}
	return nil
}

// QueryCurrentPlanRequest is the request type for the Query/AppliedPlan RPC
// method.
type QueryAppliedPlanRequest struct {
	// name is the name of the applied plan to query for.
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
}

func (m *QueryAppliedPlanRequest) Reset()         { *m = QueryAppliedPlanRequest{} }
func (m *QueryAppliedPlanRequest) String() string { return proto.CompactTextString(m) }
func (*QueryAppliedPlanRequest) ProtoMessage()    {}
func (*QueryAppliedPlanRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_4a334d07ad8374f0, []int{2}
}
func (m *QueryAppliedPlanRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryAppliedPlanRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryAppliedPlanRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryAppliedPlanRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryAppliedPlanRequest.Merge(m, src)
}
func (m *QueryAppliedPlanRequest) XXX_Size() int {
	return m.Size()
}
func (m *QueryAppliedPlanRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryAppliedPlanRequest.DiscardUnknown(m)
}

var xxx_messageInfo_QueryAppliedPlanRequest proto.InternalMessageInfo

func (m *QueryAppliedPlanRequest) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

// QueryAppliedPlanResponse is the response type for the Query/AppliedPlan RPC
// method.
type QueryAppliedPlanResponse struct {
	// height is the block height at which the plan was applied.
	Height int64 `protobuf:"varint,1,opt,name=height,proto3" json:"height,omitempty"`
}

func (m *QueryAppliedPlanResponse) Reset()         { *m = QueryAppliedPlanResponse{} }
func (m *QueryAppliedPlanResponse) String() string { return proto.CompactTextString(m) }
func (*QueryAppliedPlanResponse) ProtoMessage()    {}
func (*QueryAppliedPlanResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_4a334d07ad8374f0, []int{3}
}
func (m *QueryAppliedPlanResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryAppliedPlanResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryAppliedPlanResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryAppliedPlanResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryAppliedPlanResponse.Merge(m, src)
}
func (m *QueryAppliedPlanResponse) XXX_Size() int {
	return m.Size()
}
func (m *QueryAppliedPlanResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryAppliedPlanResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryAppliedPlanResponse proto.InternalMessageInfo

func (m *QueryAppliedPlanResponse) GetHeight() int64 {
	if m != nil {
		return m.Height
	}
	return 0
}

func init() {
	proto.RegisterType((*QueryCurrentPlanRequest)(nil), "ocea.upgrade.v1beta1.QueryCurrentPlanRequest")
	proto.RegisterType((*QueryCurrentPlanResponse)(nil), "ocea.upgrade.v1beta1.QueryCurrentPlanResponse")
	proto.RegisterType((*QueryAppliedPlanRequest)(nil), "ocea.upgrade.v1beta1.QueryAppliedPlanRequest")
	proto.RegisterType((*QueryAppliedPlanResponse)(nil), "ocea.upgrade.v1beta1.QueryAppliedPlanResponse")
}

func init() {
	proto.RegisterFile("ocea/upgrade/v1beta1/query.proto", fileDescriptor_4a334d07ad8374f0)
}

var fileDescriptor_4a334d07ad8374f0 = []byte{
	// 362 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x92, 0x3d, 0x4f, 0xfa, 0x40,
	0x1c, 0xc7, 0x39, 0xfe, 0xfc, 0x49, 0x3c, 0xb6, 0x1b, 0x10, 0x09, 0x69, 0xcc, 0x85, 0x18, 0x13,
	0xa1, 0xc7, 0xc3, 0x2b, 0x50, 0x13, 0x27, 0x07, 0x65, 0x74, 0x31, 0x07, 0xfc, 0x52, 0x1a, 0xcb,
	0xdd, 0xd1, 0xbb, 0x1a, 0x89, 0x71, 0xf1, 0x15, 0x98, 0xb8, 0xbb, 0xf9, 0x5e, 0x1c, 0x49, 0x5c,
	0x1c, 0x0d, 0xf8, 0x42, 0x0c, 0xd7, 0x62, 0x6a, 0xa0, 0xc1, 0xa9, 0x6d, 0xfa, 0x7d, 0xf8, 0xdc,
	0xb7, 0xc5, 0x74, 0x20, 0xf5, 0x58, 0x6a, 0x16, 0x29, 0x2f, 0xe4, 0x43, 0x60, 0xb7, 0xed, 0x3e,
	0x18, 0xde, 0x66, 0x93, 0x08, 0xc2, 0xa9, 0xab, 0x42, 0x69, 0x24, 0x29, 0xc7, 0x1a, 0x37, 0xd1,
	0xb8, 0x89, 0xa6, 0x5a, 0xf3, 0xa4, 0xf4, 0x02, 0x60, 0x5c, 0xf9, 0x8c, 0x0b, 0x21, 0x0d, 0x37,
	0xbe, 0x14, 0x3a, 0x76, 0x55, 0xeb, 0x19, 0xc9, 0xab, 0x14, 0xab, 0xa2, 0x7b, 0x78, 0xf7, 0x72,
	0x59, 0x75, 0x1a, 0x85, 0x21, 0x08, 0x73, 0x11, 0x70, 0xd1, 0x83, 0x49, 0x04, 0xda, 0xd0, 0x73,
	0x5c, 0x59, 0x7f, 0xa5, 0x95, 0x14, 0x1a, 0x48, 0x0b, 0x17, 0x54, 0xc0, 0x45, 0x05, 0xed, 0xa3,
	0xc3, 0x52, 0xa7, 0xe6, 0x6e, 0x26, 0x74, 0xad, 0xc7, 0x2a, 0x69, 0x33, 0x29, 0x3a, 0x56, 0x2a,
	0xf0, 0x61, 0x98, 0x2a, 0x22, 0x04, 0x17, 0x04, 0x1f, 0x83, 0x0d, 0xdb, 0xe9, 0xd9, 0x7b, 0xda,
	0x49, 0xca, 0x7f, 0xc9, 0x93, 0xf2, 0x32, 0x2e, 0x8e, 0xc0, 0xf7, 0x46, 0xc6, 0x3a, 0xfe, 0xf5,
	0x92, 0xa7, 0xce, 0x2c, 0x8f, 0xff, 0x5b, 0x13, 0x79, 0x41, 0xb8, 0x94, 0xc2, 0x26, 0x2c, 0x0b,
	0x30, 0xe3, 0xec, 0xd5, 0xd6, 0xdf, 0x0d, 0x31, 0x14, 0x6d, 0x3c, 0xbe, 0x7f, 0x3d, 0xe7, 0x0f,
	0x48, 0x9d, 0x65, 0xec, 0x3e, 0x88, 0x4d, 0xd7, 0xcb, 0x35, 0xc8, 0x2b, 0xc2, 0xa5, 0xd4, 0xd1,
	0xb6, 0x00, 0xae, 0x6f, 0xb6, 0x05, 0x70, 0xc3, 0x6a, 0xb4, 0x6b, 0x01, 0x9b, 0xe4, 0x28, 0x0b,
	0x90, 0xc7, 0x26, 0x0b, 0xc8, 0xee, 0x97, 0x5f, 0xe1, 0xe1, 0xe4, 0xec, 0x6d, 0xee, 0xa0, 0xd9,
	0xdc, 0x41, 0x9f, 0x73, 0x07, 0x3d, 0x2d, 0x9c, 0xdc, 0x6c, 0xe1, 0xe4, 0x3e, 0x16, 0x4e, 0xee,
	0xaa, 0xe1, 0xf9, 0x66, 0x14, 0xf5, 0xdd, 0x81, 0x1c, 0xaf, 0x02, 0xe3, 0x4b, 0x53, 0x0f, 0x6f,
	0xd8, 0xdd, 0x4f, 0xba, 0x99, 0x2a, 0xd0, 0xfd, 0xa2, 0xfd, 0xdb, 0xba, 0xdf, 0x01, 0x00, 0x00,
	0xff, 0xff, 0xca, 0xa9, 0x39, 0xfe, 0xef, 0x02, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// QueryClient is the client API for Query service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type QueryClient interface {
	// CurrentPlan queries the current upgrade plan.
	CurrentPlan(ctx context.Context, in *QueryCurrentPlanRequest, opts ...grpc.CallOption) (*QueryCurrentPlanResponse, error)
	// AppliedPlan queries a previously applied upgrade plan by its name.
	AppliedPlan(ctx context.Context, in *QueryAppliedPlanRequest, opts ...grpc.CallOption) (*QueryAppliedPlanResponse, error)
}

type queryClient struct {
	cc grpc1.ClientConn
}

func NewQueryClient(cc grpc1.ClientConn) QueryClient {
	return &queryClient{cc}
}

func (c *queryClient) CurrentPlan(ctx context.Context, in *QueryCurrentPlanRequest, opts ...grpc.CallOption) (*QueryCurrentPlanResponse, error) {
	out := new(QueryCurrentPlanResponse)
	err := c.cc.Invoke(ctx, "/ocea.upgrade.v1beta1.Query/CurrentPlan", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) AppliedPlan(ctx context.Context, in *QueryAppliedPlanRequest, opts ...grpc.CallOption) (*QueryAppliedPlanResponse, error) {
	out := new(QueryAppliedPlanResponse)
	err := c.cc.Invoke(ctx, "/ocea.upgrade.v1beta1.Query/AppliedPlan", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// QueryServer is the server API for Query service.
type QueryServer interface {
	// CurrentPlan queries the current upgrade plan.
	CurrentPlan(context.Context, *QueryCurrentPlanRequest) (*QueryCurrentPlanResponse, error)
	// AppliedPlan queries a previously applied upgrade plan by its name.
	AppliedPlan(context.Context, *QueryAppliedPlanRequest) (*QueryAppliedPlanResponse, error)
}

// UnimplementedQueryServer can be embedded to have forward compatible implementations.
type UnimplementedQueryServer struct {
}

func (*UnimplementedQueryServer) CurrentPlan(ctx context.Context, req *QueryCurrentPlanRequest) (*QueryCurrentPlanResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CurrentPlan not implemented")
}
func (*UnimplementedQueryServer) AppliedPlan(ctx context.Context, req *QueryAppliedPlanRequest) (*QueryAppliedPlanResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AppliedPlan not implemented")
}

func RegisterQueryServer(s grpc1.Server, srv QueryServer) {
	s.RegisterService(&_Query_serviceDesc, srv)
}

func _Query_CurrentPlan_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryCurrentPlanRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).CurrentPlan(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ocea.upgrade.v1beta1.Query/CurrentPlan",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).CurrentPlan(ctx, req.(*QueryCurrentPlanRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_AppliedPlan_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryAppliedPlanRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).AppliedPlan(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ocea.upgrade.v1beta1.Query/AppliedPlan",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).AppliedPlan(ctx, req.(*QueryAppliedPlanRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Query_serviceDesc = grpc.ServiceDesc{
	ServiceName: "ocea.upgrade.v1beta1.Query",
	HandlerType: (*QueryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CurrentPlan",
			Handler:    _Query_CurrentPlan_Handler,
		},
		{
			MethodName: "AppliedPlan",
			Handler:    _Query_AppliedPlan_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "ocea/upgrade/v1beta1/query.proto",
}

func (m *QueryCurrentPlanRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryCurrentPlanRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryCurrentPlanRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	return len(dAtA) - i, nil
}

func (m *QueryCurrentPlanResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryCurrentPlanResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryCurrentPlanResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Plan != nil {
		{
			size, err := m.Plan.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintQuery(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *QueryAppliedPlanRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryAppliedPlanRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryAppliedPlanRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Name) > 0 {
		i -= len(m.Name)
		copy(dAtA[i:], m.Name)
		i = encodeVarintQuery(dAtA, i, uint64(len(m.Name)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *QueryAppliedPlanResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryAppliedPlanResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryAppliedPlanResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Height != 0 {
		i = encodeVarintQuery(dAtA, i, uint64(m.Height))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func encodeVarintQuery(dAtA []byte, offset int, v uint64) int {
	offset -= sovQuery(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *QueryCurrentPlanRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func (m *QueryCurrentPlanResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Plan != nil {
		l = m.Plan.Size()
		n += 1 + l + sovQuery(uint64(l))
	}
	return n
}

func (m *QueryAppliedPlanRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Name)
	if l > 0 {
		n += 1 + l + sovQuery(uint64(l))
	}
	return n
}

func (m *QueryAppliedPlanResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Height != 0 {
		n += 1 + sovQuery(uint64(m.Height))
	}
	return n
}

func sovQuery(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozQuery(x uint64) (n int) {
	return sovQuery(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *QueryCurrentPlanRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
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
			return fmt.Errorf("proto: QueryCurrentPlanRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryCurrentPlanRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthQuery
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthQuery
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *QueryCurrentPlanResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
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
			return fmt.Errorf("proto: QueryCurrentPlanResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryCurrentPlanResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Plan", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
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
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Plan == nil {
				m.Plan = &Plan{}
			}
			if err := m.Plan.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthQuery
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthQuery
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *QueryAppliedPlanRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
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
			return fmt.Errorf("proto: QueryAppliedPlanRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryAppliedPlanRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Name", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
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
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Name = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthQuery
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthQuery
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *QueryAppliedPlanResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
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
			return fmt.Errorf("proto: QueryAppliedPlanResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryAppliedPlanResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Height", wireType)
			}
			m.Height = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Height |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthQuery
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthQuery
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipQuery(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowQuery
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
					return 0, ErrIntOverflowQuery
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
					return 0, ErrIntOverflowQuery
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
				return 0, ErrInvalidLengthQuery
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupQuery
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthQuery
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthQuery        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowQuery          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupQuery = fmt.Errorf("proto: unexpected end of group")
)

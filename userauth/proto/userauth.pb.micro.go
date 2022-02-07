// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: proto/userauth.proto

package userauth

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

import (
	context "context"
	api "github.com/micro/micro/v3/service/api"
	client "github.com/micro/micro/v3/service/client"
	server "github.com/micro/micro/v3/service/server"
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

// Reference imports to suppress errors if they are not otherwise used.
var _ api.Endpoint
var _ context.Context
var _ client.Option
var _ server.Option

// Api Endpoints for Userauth service

func NewUserauthEndpoints() []*api.Endpoint {
	return []*api.Endpoint{}
}

// Client API for Userauth service

type UserauthService interface {
	Call(ctx context.Context, in *Request, opts ...client.CallOption) (*Response, error)
	Stream(ctx context.Context, in *StreamingRequest, opts ...client.CallOption) (Userauth_StreamService, error)
	PingPong(ctx context.Context, opts ...client.CallOption) (Userauth_PingPongService, error)
	UserLogin(ctx context.Context, in *LoginRequest, opts ...client.CallOption) (*LoginResponse, error)
	UserLogout(ctx context.Context, in *LogoutRequest, opts ...client.CallOption) (*LogoutResponse, error)
}

type userauthService struct {
	c    client.Client
	name string
}

func NewUserauthService(name string, c client.Client) UserauthService {
	return &userauthService{
		c:    c,
		name: name,
	}
}

func (c *userauthService) Call(ctx context.Context, in *Request, opts ...client.CallOption) (*Response, error) {
	req := c.c.NewRequest(c.name, "Userauth.Call", in)
	out := new(Response)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userauthService) Stream(ctx context.Context, in *StreamingRequest, opts ...client.CallOption) (Userauth_StreamService, error) {
	req := c.c.NewRequest(c.name, "Userauth.Stream", &StreamingRequest{})
	stream, err := c.c.Stream(ctx, req, opts...)
	if err != nil {
		return nil, err
	}
	if err := stream.Send(in); err != nil {
		return nil, err
	}
	return &userauthServiceStream{stream}, nil
}

type Userauth_StreamService interface {
	Context() context.Context
	SendMsg(interface{}) error
	RecvMsg(interface{}) error
	Close() error
	Recv() (*StreamingResponse, error)
}

type userauthServiceStream struct {
	stream client.Stream
}

func (x *userauthServiceStream) Close() error {
	return x.stream.Close()
}

func (x *userauthServiceStream) Context() context.Context {
	return x.stream.Context()
}

func (x *userauthServiceStream) SendMsg(m interface{}) error {
	return x.stream.Send(m)
}

func (x *userauthServiceStream) RecvMsg(m interface{}) error {
	return x.stream.Recv(m)
}

func (x *userauthServiceStream) Recv() (*StreamingResponse, error) {
	m := new(StreamingResponse)
	err := x.stream.Recv(m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (c *userauthService) PingPong(ctx context.Context, opts ...client.CallOption) (Userauth_PingPongService, error) {
	req := c.c.NewRequest(c.name, "Userauth.PingPong", &Ping{})
	stream, err := c.c.Stream(ctx, req, opts...)
	if err != nil {
		return nil, err
	}
	return &userauthServicePingPong{stream}, nil
}

type Userauth_PingPongService interface {
	Context() context.Context
	SendMsg(interface{}) error
	RecvMsg(interface{}) error
	Close() error
	Send(*Ping) error
	Recv() (*Pong, error)
}

type userauthServicePingPong struct {
	stream client.Stream
}

func (x *userauthServicePingPong) Close() error {
	return x.stream.Close()
}

func (x *userauthServicePingPong) Context() context.Context {
	return x.stream.Context()
}

func (x *userauthServicePingPong) SendMsg(m interface{}) error {
	return x.stream.Send(m)
}

func (x *userauthServicePingPong) RecvMsg(m interface{}) error {
	return x.stream.Recv(m)
}

func (x *userauthServicePingPong) Send(m *Ping) error {
	return x.stream.Send(m)
}

func (x *userauthServicePingPong) Recv() (*Pong, error) {
	m := new(Pong)
	err := x.stream.Recv(m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (c *userauthService) UserLogin(ctx context.Context, in *LoginRequest, opts ...client.CallOption) (*LoginResponse, error) {
	req := c.c.NewRequest(c.name, "Userauth.UserLogin", in)
	out := new(LoginResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userauthService) UserLogout(ctx context.Context, in *LogoutRequest, opts ...client.CallOption) (*LogoutResponse, error) {
	req := c.c.NewRequest(c.name, "Userauth.UserLogout", in)
	out := new(LogoutResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Userauth service

type UserauthHandler interface {
	Call(context.Context, *Request, *Response) error
	Stream(context.Context, *StreamingRequest, Userauth_StreamStream) error
	PingPong(context.Context, Userauth_PingPongStream) error
	UserLogin(context.Context, *LoginRequest, *LoginResponse) error
	UserLogout(context.Context, *LogoutRequest, *LogoutResponse) error
}

func RegisterUserauthHandler(s server.Server, hdlr UserauthHandler, opts ...server.HandlerOption) error {
	type userauth interface {
		Call(ctx context.Context, in *Request, out *Response) error
		Stream(ctx context.Context, stream server.Stream) error
		PingPong(ctx context.Context, stream server.Stream) error
		UserLogin(ctx context.Context, in *LoginRequest, out *LoginResponse) error
		UserLogout(ctx context.Context, in *LogoutRequest, out *LogoutResponse) error
	}
	type Userauth struct {
		userauth
	}
	h := &userauthHandler{hdlr}
	return s.Handle(s.NewHandler(&Userauth{h}, opts...))
}

type userauthHandler struct {
	UserauthHandler
}

func (h *userauthHandler) Call(ctx context.Context, in *Request, out *Response) error {
	return h.UserauthHandler.Call(ctx, in, out)
}

func (h *userauthHandler) Stream(ctx context.Context, stream server.Stream) error {
	m := new(StreamingRequest)
	if err := stream.Recv(m); err != nil {
		return err
	}
	return h.UserauthHandler.Stream(ctx, m, &userauthStreamStream{stream})
}

type Userauth_StreamStream interface {
	Context() context.Context
	SendMsg(interface{}) error
	RecvMsg(interface{}) error
	Close() error
	Send(*StreamingResponse) error
}

type userauthStreamStream struct {
	stream server.Stream
}

func (x *userauthStreamStream) Close() error {
	return x.stream.Close()
}

func (x *userauthStreamStream) Context() context.Context {
	return x.stream.Context()
}

func (x *userauthStreamStream) SendMsg(m interface{}) error {
	return x.stream.Send(m)
}

func (x *userauthStreamStream) RecvMsg(m interface{}) error {
	return x.stream.Recv(m)
}

func (x *userauthStreamStream) Send(m *StreamingResponse) error {
	return x.stream.Send(m)
}

func (h *userauthHandler) PingPong(ctx context.Context, stream server.Stream) error {
	return h.UserauthHandler.PingPong(ctx, &userauthPingPongStream{stream})
}

type Userauth_PingPongStream interface {
	Context() context.Context
	SendMsg(interface{}) error
	RecvMsg(interface{}) error
	Close() error
	Send(*Pong) error
	Recv() (*Ping, error)
}

type userauthPingPongStream struct {
	stream server.Stream
}

func (x *userauthPingPongStream) Close() error {
	return x.stream.Close()
}

func (x *userauthPingPongStream) Context() context.Context {
	return x.stream.Context()
}

func (x *userauthPingPongStream) SendMsg(m interface{}) error {
	return x.stream.Send(m)
}

func (x *userauthPingPongStream) RecvMsg(m interface{}) error {
	return x.stream.Recv(m)
}

func (x *userauthPingPongStream) Send(m *Pong) error {
	return x.stream.Send(m)
}

func (x *userauthPingPongStream) Recv() (*Ping, error) {
	m := new(Ping)
	if err := x.stream.Recv(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (h *userauthHandler) UserLogin(ctx context.Context, in *LoginRequest, out *LoginResponse) error {
	return h.UserauthHandler.UserLogin(ctx, in, out)
}

func (h *userauthHandler) UserLogout(ctx context.Context, in *LogoutRequest, out *LogoutResponse) error {
	return h.UserauthHandler.UserLogout(ctx, in, out)
}

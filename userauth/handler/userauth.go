package handler

import (
	"context"

	log "github.com/micro/micro/v3/service/logger"

	userauth "userauth/proto"
)

type Userauth struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Userauth) Call(ctx context.Context, req *userauth.Request, rsp *userauth.Response) error {
	log.Info("Received Userauth.Call request")
	rsp.Msg = "Hello " + req.Name
	return nil
}

// Stream is a server side stream handler called via client.Stream or the generated client code
func (e *Userauth) Stream(ctx context.Context, req *userauth.StreamingRequest, stream userauth.Userauth_StreamStream) error {
	log.Infof("Received Userauth.Stream request with count: %d", req.Count)

	for i := 0; i < int(req.Count); i++ {
		log.Infof("Responding: %d", i)
		if err := stream.Send(&userauth.StreamingResponse{
			Count: int64(i),
		}); err != nil {
			return err
		}
	}

	return nil
}

// PingPong is a bidirectional stream handler called via client.Stream or the generated client code
func (e *Userauth) PingPong(ctx context.Context, stream userauth.Userauth_PingPongStream) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			return err
		}
		log.Infof("Got ping %v", req.Stroke)
		if err := stream.Send(&userauth.Pong{Stroke: req.Stroke}); err != nil {
			return err
		}
	}
}

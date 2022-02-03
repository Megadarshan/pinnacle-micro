package handler

import (
	"context"

	log "github.com/micro/micro/v3/service/logger"

	user_management "user_management/proto"
)

type User_management struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *User_management) Call(ctx context.Context, req *user_management.Request, rsp *user_management.Response) error {
	log.Info("Received User_management.Call request")
	rsp.Msg = "Hello " + req.Name
	return nil
}

// Stream is a server side stream handler called via client.Stream or the generated client code
func (e *User_management) Stream(ctx context.Context, req *user_management.StreamingRequest, stream user_management.User_management_StreamStream) error {
	log.Infof("Received User_management.Stream request with count: %d", req.Count)

	for i := 0; i < int(req.Count); i++ {
		log.Infof("Responding: %d", i)
		if err := stream.Send(&user_management.StreamingResponse{
			Count: int64(i),
		}); err != nil {
			return err
		}
	}

	return nil
}

// PingPong is a bidirectional stream handler called via client.Stream or the generated client code
func (e *User_management) PingPong(ctx context.Context, stream user_management.User_management_PingPongStream) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			return err
		}
		log.Infof("Got ping %v", req.Stroke)
		if err := stream.Send(&user_management.Pong{Stroke: req.Stroke}); err != nil {
			return err
		}
	}
}

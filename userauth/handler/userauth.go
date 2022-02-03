package handler

import (
	"context"

	errors "github.com/micro/micro/v3/service/errors"
	log "github.com/micro/micro/v3/service/logger"

	token "github.com/Megadarshan/pinnacle-micro/userauth/auth_token"
	database "github.com/Megadarshan/pinnacle-micro/userauth/database"
	userauth "github.com/Megadarshan/pinnacle-micro/userauth/proto"

	redis "github.com/Megadarshan/pinnacle-micro/redis/proto"
	"github.com/micro/micro/v3/service/client"
)

type Userauth struct {
	// httpRequest http.Request
}

// func UnaryInterceptor(
// 	ctx context.Context,
// 	req interface{},
// 	info *grpc.UnaryServerInfo,
// 	handler grpc.UnaryHandler,
// ) (interface{}, error) {
// 	log.Info("Unary Interceptor")
// 	return handler(ctx, req)
// }

// Call is a single request handler called via client.Call or the generated client code
func (e *Userauth) Call(ctx context.Context, req *userauth.Request, rsp *userauth.Response) error {
	log.Info("Received Userauth.Call request")
	rsp.Msg = "Hello " + req.Name
	return nil
}

// Call is a single request handler called via client.Call or the generated client code
func (e *Userauth) UserLogin(ctx context.Context, req *userauth.LoginRequest, rsp *userauth.LoginResponse) error {
	log.Info("Received User_auth.Call request")
	// log.Infof("%s", e.httpRequest.Header)
	var auth bool
	var er interface{}

	// var err error
	var err = database.DB.QueryRow("SELECT (password = crypt($1 , password)) as auth from users where username = $2 ",
		req.Password, req.Username).Scan(&auth)
	if err != nil {
		auth = false
	}

	rsp.LoginStatus = auth

	if auth {
		tkn, err := token.CreateToken(234, req.Username)
		if err != nil {
			println("Token not generated")
		}

		// create a new service client
		cache := redis.NewRedisService("redis", client.DefaultClient)
		// call the endpoint Redis.Set
		crsp, err := cache.Set(context.Background(), &redis.SetRequest{
			Key:   tkn.AccessUuid,
			Value: tkn.AccessToken,
			Life:  tkn.AtExpires,
		})
		log.Info(crsp)

		// rsp.Token, err = json.Marshal(token)
		if err != nil {
			println("Token not cached")
			log.Info("Token not cached")
		}
		// rsp.Token = "Username:" + req.Username + " | Password:" + req.Password
		rsp.AccessToken = tkn.AccessToken
		rsp.RefreshToken = tkn.RefreshToken
		rsp.AccessUuid = tkn.AccessUuid
		rsp.RefreshUuid = tkn.RefreshUuid
		rsp.AtExpires = tkn.AtExpires
		rsp.RtExpires = tkn.RtExpires
	} else {
		rr := errors.Unauthorized("userauth.UserLogin", "Username and Password combination not found..", er)
		return rr
	}

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

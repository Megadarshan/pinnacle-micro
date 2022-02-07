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
	"google.golang.org/grpc/metadata"
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
	var userId uint64

	var er interface{}
	//**********************************
	var values []string
	var headtoken, username string

	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		values = md.Get("authorization")
		log.Info("Header => ", md)
	}

	if len(values) > 0 {
		headtoken = values[0]
	}
	log.Info("Header -> ", headtoken)

	// // Loop over header names
	// for name, values := range md {
	// 	// Loop over all values for the name.
	// 	for _, value := range values {
	// 		log.Info("Header - (", name, " : ", value, ")")
	// 	}
	// }
	//**********************************
	// var err error
	var err = database.DB.QueryRow("SELECT id, username, (password = crypt($1 , password)) as auth from users where username = $2 ",
		req.Password, req.Username).Scan(&userId, &username, &auth)
	if err != nil {
		auth = false
	}

	rsp.LoginStatus = auth

	if auth {
		tkn, err := token.CreateToken(userId, username)
		if err != nil {
			println("Token not generated")
		}

		// create a new service client
		cache := redis.NewRedisService("redis", client.DefaultClient)
		// call the endpoint Redis.Set
		setCacheAT, err := cache.Set(context.Background(), &redis.SetRequest{
			Key:   tkn.AccessUuid,
			Value: tkn.AccessToken,
			Life:  tkn.AtExpires,
		})
		log.Info(setCacheAT)
		if err != nil {
			println("Access Token not cached")
			log.Info("Access Token not cached")
		}
		setCacheRT, err := cache.Set(context.Background(), &redis.SetRequest{
			Key:   tkn.RefreshUuid,
			Value: tkn.RefreshToken,
			Life:  tkn.RtExpires,
		})
		log.Info(setCacheRT)
		// rsp.Token, err = json.Marshal(token)
		if err != nil {
			println("Refresh Token not cached")
			log.Info("Refresh Token not cached")
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

// Call is a single request handler called via client.Call or the generated client code
func (e *Userauth) UserLogout(ctx context.Context, req *userauth.LogoutRequest, rsp *userauth.LogoutResponse) error {
	// var er interface{}
	log.Info("Received Userauth.UserLogout request")
	// create a new service client
	// err := token.TokenValidate(ctx)
	// if err != nil {
	// 	log.Info("Error : ", err.Error())
	// 	return errors.Unauthorized("userauth.UserLogout", "Unauthorized to perform Logout.. ("+err.Error()+")", er)
	// }

	cache := redis.NewRedisService("redis", client.DefaultClient)
	for i, uuid := range req.Uuids {
		log.Info(i, uuid)
		delResp, err := cache.Delete(context.Background(), &redis.DeleteRequest{Key: uuid})
		if err != nil {
			log.Info("Token not found")
		}
		log.Info(delResp)

	}

	log.Info("Successfully Logged Out...")
	rsp.Msg = "Successfully Logged Out... "
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

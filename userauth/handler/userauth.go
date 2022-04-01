package handler

import (
	"context"

	errors "github.com/micro/micro/v3/service/errors"
	log "github.com/micro/micro/v3/service/logger"

	database "github.com/Megadarshan/pinnacle-micro/userauth/database"
	userauth "github.com/Megadarshan/pinnacle-micro/userauth/proto"

	token "github.com/Megadarshan/pinnacle-micro/managetoken/proto"
	redis "github.com/Megadarshan/pinnacle-micro/redis/proto"
	"github.com/micro/micro/v3/service/client"
	"google.golang.org/grpc/metadata"
)

type Userauth struct {
	// httpRequest http.Request
}

func TokenService() token.ManagetokenService {
	service := token.NewManagetokenService("managetoken", client.DefaultClient)
	return service
}

func RedisService() redis.RedisService {
	service := redis.NewRedisService("redis", client.DefaultClient)
	return service
}

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

		token, err := TokenService().CreateToken(context.Background(), &token.CreateTokenRequest{
			UserId: int64(userId),
			// Username: username,

		})
		if err != nil {
			return errors.Unauthorized("userauth.UserLogin", "Token Not Generated", er)
		}

		// create a new service client
		// cache := redis.NewRedisService("redis", client.DefaultClient)
		setCacheAT, err := RedisService().Set(context.Background(), &redis.SetRequest{
			Key:   token.AccessUuid,
			Value: token.AccessToken,
			Life:  token.AtExpires,
		})
		// call the endpoint Redis.Set
		if err != nil {
			log.Info("Access Token ERROR: " + err.Error())
			return err
		}
		log.Info(setCacheAT)

		setCacheRT, err := RedisService().Set(context.Background(), &redis.SetRequest{
			Key:   token.RefreshUuid,
			Value: token.RefreshToken,
			Life:  token.RtExpires,
		})
		// rsp.Token, err = json.Marshal(token)
		if err != nil {
			log.Info("Refresh Token ERROR: " + err.Error())
			return err
		}
		log.Info(setCacheRT)

		rsp.AccessToken = token.AccessToken
		rsp.RefreshToken = token.RefreshToken
		rsp.AccessUuid = token.AccessUuid
		rsp.RefreshUuid = token.RefreshUuid
		rsp.AtExpires = token.AtExpires
		rsp.RtExpires = token.RtExpires
	} else {
		return errors.Unauthorized("userauth.UserLogin", "%s", "Username and Password combination not found")
	}
	return nil
}

// Call is a single request handler called via client.Call or the generated client code
func (e *Userauth) UserLogout(ctx context.Context, req *userauth.LogoutRequest, rsp *userauth.LogoutResponse) error {
	// var er interface{}
	log.Info("Received Userauth.UserLogout request")

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

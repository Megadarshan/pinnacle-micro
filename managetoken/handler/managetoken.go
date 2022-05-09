package handler

import (
	"encoding/json"
	"strings"

	managetoken "github.com/Megadarshan/pinnacle-micro/managetoken/proto"

	"context"
	"time"

	redis "github.com/Megadarshan/pinnacle-micro/redis/proto"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/micro/micro/v3/service/client"
	errors "github.com/micro/micro/v3/service/errors"
	log "github.com/micro/micro/v3/service/logger"
	"github.com/twinj/uuid"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

func RedisService() redis.RedisService {
	service := redis.NewRedisService("redis", client.DefaultClient)
	return service
}

const AccessTokenSign = "126F25455F5ABB6CF1DAB9EA24C77"
const RefreshTokenSign = "A2FD9D367EB21BAE3EEEA8B79F9EE"

type Managetoken struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Managetoken) CreateToken(ctx context.Context, req *managetoken.CreateTokenRequest, rsp *managetoken.CreateTokenResponse) error {
	log.Info("Received Managetoken.CreateToken request")
	var dummyInterface interface{}

	md, ok := metadata.FromIncomingContext(ctx)

	if !ok {
		return errors.Unauthorized("managetoken.ValidateToken", "%s", "No Context found")
	}

	for i, j := range md {
		log.Info(i, " - ", j)
	}

	AtExpires := time.Now().Add(time.Minute * 2).Unix()
	AccessUuid := uuid.NewV4().String()

	RtExpires := time.Now().Add(time.Minute * 5).Unix()
	RefreshUuid := uuid.NewV4().String()

	if req.StaySignedIn {
		RtExpires = time.Now().AddDate(0, 0, 30).Unix()
	}

	var err error

	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_uuid"] = AccessUuid
	atClaims["exp"] = AtExpires
	atClaims["user_id"] = req.UserId
	atClaims["username"] = req.Username
	// atClaims["struct"] = req.Claims

	var inInterface map[string]interface{}
	inrec, _ := json.Marshal(req.Claims)
	json.Unmarshal(inrec, &inInterface)

	// iterate through inrecs
	for field, val := range inInterface {
		atClaims[field] = val
	}

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	AccessToken, err := at.SignedString([]byte(AccessTokenSign))
	if err != nil {
		return errors.InternalServerError("Managetoken.CreateToken", "Access token not generated"+err.Error(), dummyInterface)
	}

	//Creating Refresh Token
	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = RefreshUuid
	rtClaims["user_id"] = req.UserId
	rtClaims["username"] = req.Username
	rtClaims["exp"] = RtExpires
	// rtClaims["struct"] = req.Claims
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	RefreshToken, err := rt.SignedString([]byte(RefreshTokenSign))
	if err != nil {
		return errors.InternalServerError("Managetoken.CreateToken", "Refresh token not generated"+err.Error(), dummyInterface)
	}

	rsp.AccessUuid = AccessUuid
	rsp.AccessToken = AccessToken
	rsp.AtExpires = AtExpires
	rsp.RefreshUuid = RefreshUuid
	rsp.RefreshToken = RefreshToken
	rsp.RtExpires = RtExpires

	return nil
}

func (e *Managetoken) ValidateToken(ctx context.Context, req *managetoken.ValidateTokenRequest, rsp *managetoken.ValidateTokenResponse) error {
	log.Info("Received Managetoken.ValidateToken request\n--------------------------------------------")

	md, ok := metadata.FromIncomingContext(ctx)

	if !ok {
		return errors.Unauthorized("managetoken.ValidateToken", "%s", "No Context found")
	}

	for i, j := range md {
		log.Info(i, " - ", j)
	}

	accessuuid := md["accessuuid"]
	if len(accessuuid) != 1 {
		return errors.Unauthorized("managetoken.ValidateToken", "%s", "No AuthToken ID found in request")
	}

	AccessToken, err := RedisService().Get(context.Background(), &redis.GetRequest{
		Key: accessuuid[0],
	})

	if err != nil {
		return errors.Unauthorized("managetoken.ValidateToken", "%s", "Unable to get stored token")
	}

	if AccessToken.Value == "" {
		return errors.Unauthorized("managetoken.ValidateToken", "%s", "Access token not available or expired")
	}

	AccessTokenIP, err := RedisService().Get(context.Background(), &redis.GetRequest{
		Key: accessuuid[0] + "-IP",
	})
	if err != nil {
		log.Info("Unable to get stored IP")
	}

	if strings.Split(md.Get(":authority")[0], ":")[0] != AccessTokenIP.Value {
		return errors.Unauthorized("managetoken.ValidateToken", "%s", "Requested detected from a new device/browser. Please login again")
	}

	if len(md["authorization"]) != 1 {
		return errors.Unauthorized("managetoken.ValidateToken", "%s", "No AuthToken found")
	}

	if strings.Split(md["authorization"][0], " ")[1] != AccessToken.Value {
		return errors.Unauthorized("managetoken.ValidateToken", "%s", "Access token seems tampered")
	}

	auths, ok := md[":authority"]
	if !ok {
		// return nil, status.Error(codes.InvalidArgument, "no authority header")
		return errors.Unauthorized("managetoken.ValidateToken", "%s", "No Authority header")
	}
	if len(auths) != 1 {
		// return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("no authority header, auths = %v", auths))
		return errors.Unauthorized("managetoken.ValidateToken", "no authority header, auths = %v", auths)
	}
	// if auths[0] != expectedAuthority {
	// 	// return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("invalid authority header %v, expected %v", auths[0], expectedAuthority))
	// 	return errors.Unauthorized("managetoken.ValidateToken", "%s", "invalid authority header *, expected *")
	// }

	p, _ := peer.FromContext(ctx)
	req_ip := p.Addr.String()
	log.Info("Request from IP : ", req_ip)

	// if ok {
	// 	values := md.Get("authorization")
	// 	// if len(values) > 0{
	// 	// 	return errors.BadRequest("userauth.UserLogin", "Authorization header not found in the request..", er)
	// 	// }
	// } else {
	// 	return errors.BadRequest("userauth.UserLogin", "Header not found in the request..", er)
	// }

	return nil
}

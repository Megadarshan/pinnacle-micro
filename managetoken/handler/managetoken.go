package handler

import (
	managetoken "github.com/Megadarshan/pinnacle-micro/managetoken/proto"

	"context"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	errors "github.com/micro/micro/v3/service/errors"
	log "github.com/micro/micro/v3/service/logger"
	"github.com/twinj/uuid"
)

const AccessTokenSign = "126F25455F5ABB6CF1DAB9EA24C77"
const RefreshTokenSign = "A2FD9D367EB21BAE3EEEA8B79F9EE"

type Managetoken struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Managetoken) CreateToken(ctx context.Context, req *managetoken.CreateTokenRequest, rsp *managetoken.CreateTokenResponse) error {
	log.Info("Received Managetoken.CreateToken request")
	var dummyInterface interface{}
	AtExpires := time.Now().Add(time.Minute * 5).Unix()
	AccessUuid := uuid.NewV4().String()

	RtExpires := time.Now().Add(time.Minute * 30).Unix()
	RefreshUuid := uuid.NewV4().String()

	var err error

	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_uuid"] = AccessUuid
	atClaims["exp"] = AtExpires
	atClaims["user_id"] = req.UserId
	atClaims["username"] = req.Username
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	AccessToken, err := at.SignedString([]byte(AccessTokenSign))
	if err != nil {
		return errors.BadRequest("Managetoken.CreateToken", "Access token not generated"+err.Error(), dummyInterface)
	}

	//Creating Refresh Token
	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = RefreshUuid
	rtClaims["user_id"] = req.UserId
	rtClaims["username"] = req.Username
	rtClaims["exp"] = RtExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	RefreshToken, err := rt.SignedString([]byte(RefreshTokenSign))
	if err != nil {
		return errors.BadRequest("Managetoken.CreateToken", "Refresh token not generated"+err.Error(), dummyInterface)
	}

	rsp.AccessUuid = AccessUuid
	rsp.AccessToken = AccessToken
	rsp.AtExpires = AtExpires
	rsp.RefreshUuid = RefreshUuid
	rsp.RefreshToken = RefreshToken
	rsp.RtExpires = RtExpires

	return nil
}

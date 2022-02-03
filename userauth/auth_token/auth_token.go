package auth_token

import (
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/twinj/uuid"
)

const (
	access_secret = "ACCESS_SECRET"
)

type AccessDetails struct {
	AccessUuid string
	UserId     uint64
}

// var Client *redis.Client

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AccessUuid   string
	RefreshUuid  string
	AtExpires    int64
	RtExpires    int64
}

// EncodeJWT serialize data into a jwt token using a secret
// This secret must match with the client's secret who's generating the token

func CreateToken(userid uint64, username string) (*TokenDetails, error) {
	td := &TokenDetails{}
	td.AtExpires = time.Now().Add(time.Minute * 15).Unix()
	td.AccessUuid = uuid.NewV4().String()

	td.RtExpires = time.Now().Add(time.Hour * 24).Unix()
	td.RefreshUuid = uuid.NewV4().String()

	var err error
	//Creating Access Token
	// os.Setenv("ACCESS_SECRET", "jdnfksdmfksd") //this should be in an env file
	// Err := godotenv.Load(".env")
	// if Err != nil {
	// 	log.Fatalf("Error loading .env file")
	// }
	// sccret_key, ok := os.LookupEnv(access_secret)
	// if !ok {
	// 	panic("ACCESS_SECRET not found in env")
	// }
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_uuid"] = td.AccessUuid
	atClaims["exp"] = td.AtExpires
	atClaims["user_id"] = userid
	atClaims["username"] = username
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte("mcmvmkmsdnfsdmfdsjf"))
	if err != nil {
		return nil, err
	}
	//Creating Refresh Token
	// os.Setenv("REFRESH_SECRET", "mcmvmkmsdnfsdmfdsjf") //this should be in an env file
	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = td.RefreshUuid
	rtClaims["user_id"] = userid
	rtClaims["exp"] = td.RtExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(os.Getenv("REFRESH_SECRET")))
	if err != nil {
		return nil, err
	}
	return td, nil
}

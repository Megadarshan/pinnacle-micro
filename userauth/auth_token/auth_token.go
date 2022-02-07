package auth_token

import (
	"context"
	"os"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	errors "github.com/micro/micro/v3/service/errors"
	log "github.com/micro/micro/v3/service/logger"
	"github.com/twinj/uuid"
	"google.golang.org/grpc/metadata"
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
	td.AtExpires = time.Now().Add(time.Minute * 5).Unix()
	td.AccessUuid = uuid.NewV4().String()

	td.RtExpires = time.Now().Add(time.Minute * 30).Unix()
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

func TokenValidate(ctx context.Context) error {
	var values []string
	var headtoken string
	var er interface{}

	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		values = md.Get("authorization")
		// if len(values) > 0{
		// 	return errors.BadRequest("userauth.UserLogin", "Authorization header not found in the request..", er)
		// }
	} else {
		return errors.BadRequest("userauth.UserLogin", "Header not found in the request..", er)
	}

	if len(values) > 0 {
		headtoken = values[0]
		strArr := strings.Split(headtoken, " ")
		if len(strArr) == 2 {
			headtoken = strArr[1]
		} else {
			return errors.BadRequest("userauth.UserLogin", "Bearer Token not found in the request..", er)
		}
	} else {
		return errors.BadRequest("userauth.UserLogin", "Authorization header not found in the request..", er)
	}

	token, err := jwt.Parse(headtoken, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			msg := "Unexpected signing method: "
			log.Info(msg)
			return nil, errors.BadRequest("userauth.UserLogin", msg, token.Header["alg"])
		}
		// td.AccessToken, err = at.SignedString([]byte("mcmvmkmsdnfsdmfdsjf"))
		//return []byte(os.Getenv("ACCESS_SECRET")), nil
		return []byte("mcmvmkmsdnfsdmfdsjf"), nil
	})
	if err != nil {
		return err
	}

	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return err
	}

	return nil
}

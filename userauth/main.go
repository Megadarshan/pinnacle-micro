package main

import (
	"github.com/Megadarshan/pinnacle-micro/userauth/database"
	"github.com/Megadarshan/pinnacle-micro/userauth/handler"
	pb "github.com/Megadarshan/pinnacle-micro/userauth/proto"

	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"
)

func main() {

	database.InitDb()
	// Create service
	srv := service.New(
		service.Name("userauth"),
		service.Version("latest"),
	)

	// Register handler
	pb.RegisterUserauthHandler(srv.Server(), new(handler.Userauth))

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}

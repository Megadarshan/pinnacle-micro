package main

import (
	"user_management/handler"
	pb "user_management/proto"

	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("user_management"),
		service.Version("latest"),
	)

	// Register handler
	pb.RegisterUser_managementHandler(srv.Server(), new(handler.User_management))

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}

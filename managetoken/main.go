package main

import (
	"github.com/Megadarshan/pinnacle-micro/managetoken/handler"
	pb "github.com/Megadarshan/pinnacle-micro/managetoken/proto"

	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("managetoken"),
		service.Version("latest"),
	)

	// Register handler
	pb.RegisterManagetokenHandler(srv.Server(), new(handler.Managetoken))

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}

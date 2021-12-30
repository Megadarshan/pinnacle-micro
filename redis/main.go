package main

import (
	"redis/handler"
	pb "redis/proto"

	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("redis"),
		service.Version("latest"),
	)

	// Register handler
	pb.RegisterRedisHandler(srv.Server(), new(handler.Redis))

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}

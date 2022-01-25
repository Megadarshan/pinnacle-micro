package main

import (
	"github.com/Megadarshan/pinnacle-micro/redis/handler"
	pb "github.com/Megadarshan/pinnacle-micro/redis/proto"
	"github.com/Megadarshan/pinnacle-micro/redis/redis_ops"

	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"
)

func main() {

	redis_ops.InitRedis()
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

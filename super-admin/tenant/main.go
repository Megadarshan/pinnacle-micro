package main

import (
	"github.com/Megadarshan/pinnacle-micro/super-admin/tenant/database"
	"github.com/Megadarshan/pinnacle-micro/super-admin/tenant/handler"
	pb "github.com/Megadarshan/pinnacle-micro/super-admin/tenant/proto"

	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"
)

func main() {
	database.InitDb()
	// Create service
	srv := service.New(
		service.Name("tenant"),
		service.Version("latest"),
	)

	// Register handler
	pb.RegisterTenantHandler(srv.Server(), new(handler.Tenant))

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}

package main

import (
	"context"
	"log"
	"net"
	"testing"

	"github.com/Megadarshan/pinnacle-micro/super-admin/tenant/handler"
	pb "github.com/Megadarshan/pinnacle-micro/super-admin/tenant/proto"
	"github.com/micro/micro/v3/service"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func init() {
	lis = bufconn.Listen(bufSize)
	s := service.New(
		service.Name("tenant"),
		service.Version("latest"),
	)
	// pb.RegisterGreeterServer(s, &server{})
	pb.RegisterTenantHandler(s.Server(), new(handler.Tenant))
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func TestSayHello(t *testing.T) {
	ctx := context.Background()
	conn, err := service.DialContext(ctx, "bufnet", service.WithContextDialer(bufDialer), service.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()
	client := pb.NewTenantService(conn)
	resp, err := client.GetTenantStatus(ctx, &pb.TenantStatusRequest{})
	if err != nil {
		t.Fatalf("SayHello failed: %v", err)
	}
	log.Printf("Response: %+v", resp)
	// Test for output here.
}

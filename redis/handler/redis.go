package handler

import (
	"context"

	"github.com/Megadarshan/pinnacle-micro/redis/redis_ops"

	log "github.com/micro/micro/v3/service/logger"

	redis "github.com/Megadarshan/pinnacle-micro/redis/proto"
)

type Redis struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Redis) Set(ctx context.Context, req *redis.SetRequest, rsp *redis.SetResponse) error {
	log.Info("Received Redis.Set request")
	err := redis_ops.Set(req.Key, req.Value, req.Life)
	if err != nil {
		rsp.Msg = "Something went wrong"
	}
	rsp.Msg = "Set -> {Key:" + req.Key + " Value:" + req.Value + "} Complete"
	// redis_ops.InitRedis()
	return nil
}

// Call is a single request handler called via client.Call or the generated client code
func (e *Redis) Get(ctx context.Context, req *redis.GetRequest, rsp *redis.GetResponse) error {
	log.Info("Received Redis.Get request")
	rsp.Value = "Hello " + req.Key
	return nil
}

// Call is a single request handler called via client.Call or the generated client code
func (e *Redis) Delete(ctx context.Context, req *redis.DeleteRequest, rsp *redis.DeleteResponse) error {
	log.Info("Received Redis.Delete request")
	rsp.Msg = "Hello " + req.Key
	return nil
}

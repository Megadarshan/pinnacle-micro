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
	resp, err := redis_ops.Get(req.Key)

	if err != nil {
		log.Info("Something went wrong")
	}
	rsp.Value = resp
	return nil
}

// Call is a single request handler called via client.Call or the generated client code
func (e *Redis) Delete(ctx context.Context, req *redis.DeleteRequest, rsp *redis.DeleteResponse) error {
	log.Info("Received Redis.Delete request")
	err := redis_ops.Remove(req.Key)
	if err != nil {
		log.Info("Value Not Removed")
		rsp.Msg = "Value Not Removed"
		return err
	}
	rsp.Msg = "Value Removed"
	return nil
}

func (e *Redis) Exists(ctx context.Context, req *redis.ExistsRequest, rsp *redis.ExistsResponse) error {
	log.Info("Received Redis.Exists request")
	exists, err := redis_ops.Exists(req.Key)
	if err != nil {
		log.Info(err.Error())
		rsp.Exists = exists
		return err
	}
	rsp.Exists = exists
	return nil
}

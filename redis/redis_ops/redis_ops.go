package redis_ops

import (
	"fmt"
	"log"
	"os"
	"time"

	// redis "github.com/go-redis/redis/v8"
	redis "github.com/gomodule/redigo/redis"
)

var pool *redis.Pool

func InitRedis() {
	print("Init Redis")
	// init redis connection pool
	initPool()

	// bootstramp some data to redis
	initStore()
}

func initPool() {
	pool = &redis.Pool{
		MaxIdle:   80,
		MaxActive: 12000,
		Dial: func() (redis.Conn, error) {
			// You can use redis.DialOption: redis.DialDatabase, redis.DialPassword within redis.Dial()
			conn, err := redis.Dial("tcp", "localhost:6379", redis.DialDatabase(1))
			if err != nil {
				log.Printf("ERROR: fail init redis: %s", err.Error())
				os.Exit(1)
			}
			return conn, err
		},
	}
}

func initStore() {
	// get conn and put back when exit from method
	conn := pool.Get()
	defer conn.Close()

	// macs := []string{"00000C  Cisco", "00000D  FIBRONICS", "00000E  Fujitsu",
	// 	"00000F  Next", "000010  Hughes"}
	// for _, mac := range macs {
	// 	pair := strings.Split(mac, "  ")
	// 	Set(pair[0], pair[1])
	// }
}

func Ping(conn redis.Conn) {
	_, err := redis.String(conn.Do("PING"))
	if err != nil {
		log.Printf("ERROR: fail ping redis conn: %s", err.Error())
		os.Exit(1)
	}
}

func Set(key string, val string, expiration int64) error {
	// get conn and put back when exit from method
	conn := pool.Get()
	defer conn.Close()

	at := time.Unix(expiration, 0)
	now := time.Now()
	expSeconds := int(at.Sub(now).Seconds())

	_, err := conn.Do("SET", key, val, "EX", expSeconds)

	if err != nil {
		log.Printf("ERROR: fail set key %s, val %s, error %s", key, val, err.Error())
		return err
	}

	return nil
}

func Get(key string) (string, error) {
	// get conn and put back when exit from method
	conn := pool.Get()
	defer conn.Close()

	value, err := redis.String(conn.Do("GET", key))

	if err == redis.ErrNil {
		log.Printf("%s : Alert! this Key does not exist\n", key)
		return ("Alert! this Key does not exist\n"), err
	} else if err != nil {
		log.Printf("ERROR: fail get key %s, error %s", key, err.Error())
		return ("An Unknown error encountered"), err
	} else {
		return value, err
	}
}

func Remove(key string) error {
	// get conn and put back when exit from method
	conn := pool.Get()
	defer conn.Close()

	_, err := conn.Do("DEL", key)
	return err
}

func Exists(key string) (bool, error) {

	conn := pool.Get()
	defer conn.Close()

	ok, err := redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		return ok, fmt.Errorf("error checking if key %s exists: %v", key, err)
	}
	return ok, err
}

func Sadd(key string, val string) error {
	// get conn and put back when exit from method
	conn := pool.Get()
	defer conn.Close()

	_, err := conn.Do("SADD", key, val)
	if err != nil {
		log.Printf("ERROR: fail add val %s to set %s, error %s", val, key, err.Error())
		return err
	}

	return nil
}

func Smembers(key string) ([]string, error) {
	// get conn and put back when exit from method
	conn := pool.Get()
	defer conn.Close()

	s, err := redis.Strings(conn.Do("SMEMBERS", key))
	if err != nil {
		log.Printf("ERROR: fail get set %s , error %s", key, err.Error())
		return nil, err
	}

	return s, nil
}

// func main() {
// 	// initialize redis pool and bootstrap redis
// 	initRedis()

// 	// get value which exists
// 	log.Printf(get("00000E"))

// 	// get value which does not exists
// 	log.Printf(get("0000E"))

// 	// add members to set
// 	sadd("mystiko", "0000E")
// 	sadd("mystiko", "0000D")

// 	// get memebers of set
// 	s, _ := smembers("mystiko")
// 	log.Printf("%v", s)
// }

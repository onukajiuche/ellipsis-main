package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"brief/internal/config"
	"brief/internal/constant"
	"brief/utility"

	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
)

var (
	Rds *redis.Client
	Ctx = context.Background()
)

func SetupRedis() {
	logger := log.New()
	getConfig := config.GetConfig()
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%v:%v", getConfig.RedisHost, getConfig.RedisPort),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	if err := rdb.Ping(Ctx).Err(); err != nil {
		logger.Printf("%v:%v", getConfig.RedisHost, getConfig.RedisPort)
		logger.Fatalln("Redis db error: ", err)
	}
	pong, _ := rdb.Ping(Ctx).Result()
	logger.Println("Redis says: ", pong)
	Rds = rdb
	logger.Info("Redis CONNECTION ESTABLISHED")

	// Add a counter to redis to aid hashing algorithm
	rd := GetRedisDb()
	cnt, err := rd.RedisGet(constant.CounterKey)
	if err != nil {
		if errors.Is(err, redis.ErrClosed) {
			logger.Fatal(err)
		} else {
			logger.Info("Redis COUNTER VARIABLE NOT SET")
		}
		return
	}

	var counter int64
	json.Unmarshal(cnt, &counter)
	utility.Counter = counter
}

// StoreCounter stores the current value of the counter variable in redis
func StoreCounter() {
	logger := log.New()
	rd := GetRedisDb()
	err := rd.RedisSet(constant.CounterKey, utility.Counter)
	if err != nil {
		logger.Error("Redis COULD NOT STORE COUNTER")
		return
	}
	logger.Info("Redis SUCCESSFULLY STORED REDIS COUNTER")
}

type Redis struct {
	Rdb *redis.Client
}

func GetRedisDb() *Redis {
	return &Redis{Rdb: Rds}
}

func (rdb *Redis) RedisSet(key string, value interface{}) error {
	serialized, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return rdb.Rdb.Set(Ctx, key, serialized, 24*time.Hour).Err()
}

func (rdb *Redis) RedisGet(key string) ([]byte, error) {
	serialized, err := rdb.Rdb.Get(Ctx, key).Bytes()
	return serialized, err
}

func (rdb *Redis) RedisDelete(key string) (int64, error) {
	deleted, err := rdb.Rdb.Del(Ctx, key).Result()
	if err != nil {
		return 0, err
	}
	return deleted, nil
}

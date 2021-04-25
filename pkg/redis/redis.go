package redis

import (
	"context"
	_ "crypto/tls"
	"github.com/go-redis/redis/v8"
	"github.com/smiletrl/micro_ecommerce/pkg/config"
	_ "github.com/smiletrl/micro_ecommerce/pkg/constants"
	"time"
)

// DB creates a new redis client
func DB(cfg config.Config, position int) *redis.Client {
	// redis service
	redisOptions := &redis.Options{
		Addr:     cfg.Redis.Endpoint + ":" + cfg.Redis.Port,
		Password: cfg.Redis.Password,
		DB:       position,
	}
	//stage := cfg.Stage
	//if stage != constants.StageLocal && stage != constants.StageGithub && stage != constants.StageK8s {
	//redisOptions.TLSConfig = &tls.Config{
	//	InsecureSkipVerify: false,
	//}
	//}

	rdb := redis.NewClient(redisOptions)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if _, err := rdb.Ping(ctx).Result(); err != nil {
		panic(err)
	}
	return rdb
}

// New creates the first/default redis database client
func New(cfg config.Config) *redis.Client {
	return DB(cfg, 0)
}

// Test creates the second redis database client
func Test(cfg config.Config) *redis.Client {
	return DB(cfg, 1)
}
package redis

import (
	"context"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
)

func GetNormalConnection(ctx context.Context, host, password string, database, timeout, poolSize int) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:        host,
		Password:    password,
		DB:          database,
		ReadTimeout: time.Duration(timeout) * time.Second, // This also sets the write timeout
		PoolSize:    poolSize,
	})
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return client, err
}

func GetClusterConnection(ctx context.Context, host string, timeout int) (*redis.ClusterClient, error) {
	client := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: strings.Split(host, ","),
	})
	_, err := client.Ping(ctx).Result()

	return client, err
}

func HealthCheck(ctx context.Context, client redis.Cmdable) error {
	_, err := client.Ping(ctx).Result()
	return err
}

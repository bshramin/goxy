package redis

import (
	"context"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

func GetNormalConnection(ctx context.Context, host, password string, database, timeout, poolSize int) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:        host,
		Password:    password,
		DB:          database,
		ReadTimeout: time.Duration(timeout) * time.Second, // This also sets the write timeout
		PoolSize:    poolSize,
	})
	_, err := client.Ping(ctx).Result()
	if err != nil {
		logrus.Errorf("ping to normal redis failed, host: %s", host)
		panic(err)
	}
	logrus.Infof("successfully connected to redis-normal host=%s", host)
	return client
}

func GetClusterConnection(ctx context.Context, host string, timeout int) *redis.ClusterClient {
	client := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: strings.Split(host, ","),
	})
	_, err := client.Ping(ctx).Result()
	if err != nil {
		logrus.Errorf("ping to cluster redis failed, host: %s", host)
		panic(err)
	}
	logrus.Infof("successfully connected to redis-cluster host=%s", host)
	return client
}

func HealthCheck(ctx context.Context, client redis.Cmdable) error {
	_, err := client.Ping(ctx).Result()
	return err
}

package redis

import (
	"strings"
	"time"

	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
)

func GetNormalConnection(host, password string, database, timeout, poolSize int) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:        host,
		Password:    password,
		DB:          database,
		ReadTimeout: time.Duration(timeout) * time.Second, // This also sets the write timeout
		PoolSize:    poolSize,
	})
	_, err := client.Ping().Result()
	if err != nil {
		logrus.Errorf("ping to normal redis failed, host: %s", host)
		panic(err)
	}
	logrus.Infof("successfully connected to redis-normal host=%s", host)
	return client
}

func GetClusterConnection(host string, timeout int) *redis.ClusterClient {
	client := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: strings.Split(host, ","),
	})
	_, err := client.Ping().Result()
	if err != nil {
		logrus.Errorf("ping to cluster redis failed, host: %s", host)
		panic(err)
	}
	logrus.Infof("successfully connected to redis-cluster host=%s", host)
	return client
}

func HealthCheck(client redis.Cmdable) error {
	_, err := client.Ping().Result()
	return err
}

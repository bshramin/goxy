package redis

import (
	"strings"
	"time"

	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
)

func GetNormalConnection(name, host string, port, database, timeout, poolSize int, password, kind string) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:        host,
		Password:    password,
		DB:          database,
		ReadTimeout: time.Duration(timeout) * time.Second, // This also sets the write timeout
		PoolSize:    poolSize,
	})
	_, err := client.Ping().Result()
	if err != nil {
		logrus.Errorf("ping to normal redis failed to %s, host: %s", name, host)
		return client, err
	}
	logrus.Infof("successfully connected to redis-%s host=%s", kind, host)
	return client, nil
}

func GetClusterConnection(name, host string, timeout int, kind string) (*redis.ClusterClient, error) {
	client := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: strings.Split(host, ","),
	})
	_, err := client.Ping().Result()
	if err != nil {
		logrus.Errorf("ping to cluster redis failed to %s, host: %s", name, host)
		return client, err
	}
	logrus.Infof("successfully connected to redis-%s host=%s", kind, host)
	return client, nil
}

func HealthCheck(client redis.Cmdable) error {
	_, err := client.Ping().Result()
	return err
}

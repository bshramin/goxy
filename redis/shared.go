package redis

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/bshramin/goxy"
	"github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
	"github.com/sirupsen/logrus"
)

var keysList = make(map[*redis.Client]map[string]*redsync.Mutex)

// this is used when multi instances of a code try to use shared cache on redis
// and result fetch function has too much load expensive
// it will return the same result for all of them
// it will checks if the value exists it will return it
// if not it will check if any one else is trying to fetch value
// if it was the case it will not try to fetch
// but if no one of instances has tried to fetch the result it start to fetch and set a waitKey
// too inform others
func SharedFetch[K any](ctx context.Context, client *redis.Client, key string, t, retryT time.Duration, retry int, f func() (K, error)) (K, error) {
	var res K
	redisRes := client.Get(ctx, key)
	go func() { _ = Fetch(ctx, client, key, t, retryT, retry, f) }()
	err := redisRes.Err()
	if err != nil {
		return res, err
	}
	resStr, err := redisRes.Result()
	if err != nil {
		return res, err
	}
	res, err = goxy.Decode[K](resStr)
	if err != nil {
		return res, err
	}
	return res, nil
}

func Fetch[K any](ctx context.Context, client *redis.Client, key string, t, retryT time.Duration, retry int, f func() (K, error)) error {
	lock, err := CheckLock(client, key)
	if err != nil {
		logrus.Infof("goxy:SharedFetch:lock:%s", err)
		return err
	}
	opt := redsync.WithExpiry(retryT)
	opt.Apply(lock)
	if err := lock.Lock(); err != nil {
		logrus.Infof("goxy:SharedFetch:lock:%s", err)
		return err
	}
	var res K
	for i := 0; i < retry; i++ {
		res, err = f()
		if err == nil {
			break
		}
	}
	if err != nil {
		_, lErr := lock.Unlock()
		if lErr != nil {
			logrus.Errorf("goxy:SharedFetch:unlock:%s", lErr)
		}
		logrus.Errorf("goxy:SharedFetch:fetch(after %d times):%s:%s", retry, key, err.Error())
		return err
	}
	redisSetVal, err := goxy.Encode(res)
	if err != nil {
		logrus.Errorf("goxy:SharedFetch:encode:%s", err.Error())
		return err
	}
	err = client.Set(ctx, key, redisSetVal, t).Err()
	if err != nil {
		logrus.Errorf("goxy:SharedFetch:Set:%s:%s", key, err.Error())
		return err
	}
	logrus.Infof("goxy:SharedFetch:Set:%s:Successfully", key)
	return nil
}

func CheckLock(client *redis.Client, key string) (*redsync.Mutex, error) {
	lock, ok := keysList[client][key]
	if !ok {
		_, ok := keysList[client]
		if !ok {
			keysList[client] = make(map[string]*redsync.Mutex)
		}
		pool := goredis.NewPool(client)
		rs := redsync.New(pool)
		lock = rs.NewMutex(fmt.Sprintf("%s:wait", key))
		keysList[client][key] = lock
	}
	if lock.Until().After(time.Now()) {
		return nil, errors.New("lock in on")
	}
	return lock, nil
}

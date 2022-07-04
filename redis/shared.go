package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/bshramin/goxy"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

type dataFetcher[K any] func(ctx context.Context) (K, error)

// this is used when multi instances of a code try to use shared cache on redis
// and result fetch function has too much load expensive
// it will return the same result for all of them
// it will checks if the value exists it will return it
// if not it will check if any one else is trying to fetch value
// if it was the case it will not try to fetch
// but if no one of instances has tried to fetch the result it start to fetch and set a waitKey
// too inform others
func SharedFetch[K any](ctx context.Context, client redis.Cmdable, key string, t, retryT time.Duration, retry int, f dataFetcher[K]) (K, error) {
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

func Fetch[K any](ctx context.Context, client redis.Cmdable, key string, t, retryT time.Duration, retry int, f dataFetcher[K]) error {
	err := lock(ctx, client, key, retryT)
	if err != nil {
		logrus.Infof("goxy:SharedFetch:%s-lock:%v", key, err)
		return err
	}
	res, err := fetchData(ctx, f, retry)
	if err != nil {
		if e := unlock(ctx, client, key, retryT); e != nil {
			logrus.Errorf("goxy:SharedFetch:%s-unlock:%s", key, e)
		}
		logrus.Errorf("goxy:SharedFetch:fetch(after %d times):%s:%s", retry, key, err.Error())
		return err
	}
	err = setInredis(ctx, client, key, t, res)
	if err != nil {
		logrus.Error(err)
		return err
	}
	logrus.Infof("goxy:SharedFetch:Set:%s:Successfully", key)
	return nil
}

func lock(ctx context.Context, client redis.Cmdable, key string, t time.Duration) error {
	return changeLockstatus(ctx, client, key, t, true)
}

func unlock(ctx context.Context, client redis.Cmdable, key string, t time.Duration) error {
	return changeLockstatus(ctx, client, key, t, false)
}

func changeLockstatus(ctx context.Context, client redis.Cmdable, key string, t time.Duration, status bool) error {
	k := lockingKey(key)
	cmd := client.SetNX(ctx, k, status, t)
	if val, err := cmd.Result(); err != nil || val == false {
		return err
	}
	return nil
}

func fetchData[K any](ctx context.Context, fun func(ctx context.Context) (K, error), retry int) (res K, err error) {
	for i := 0; i < retry; i++ {
		res, err = fun(ctx)
		if err != nil {
			continue
		}
		return
	}
	return
}

func lockingKey(k string) string {
	return fmt.Sprintf("%s:wait", k)
}

func setInredis[K any](ctx context.Context, client redis.Cmdable, key string, t time.Duration, data K) error {
	redisSetVal, err := goxy.Encode(data)
	if err != nil {
		err = fmt.Errorf("goxy:SharedFetch:%s-encode:%s", key, err.Error())
		return err
	}
	err = client.Set(ctx, key, redisSetVal, t).Err()
	if err != nil {
		err = fmt.Errorf("goxy:SharedFetch:Set:%s:%s", key, err.Error())
		return err
	}
	return nil
}

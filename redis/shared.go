package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/bshramin/goxy"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

// this is used when multi instances of a code try to use shared cache on redis
// and result fetch function has too much load expensive
// it will return the same result for all of them
// it will checks if the value exists it will return it
// if not it will check if any one else is trying to fetch value
// if it was the case it will not try to fetch
// but if no one of instances has tried to fetch the result it start to fetch and set a waitKey
// too inform others
func SharedFetch[K any](ctx context.Context, client redis.Cmdable, key string, t, retryT time.Duration, retry int, f func(ctx context.Context) (K, error)) (K, error) {
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

func Fetch[K any](ctx context.Context, client redis.Cmdable, key string, t, retryT time.Duration, retry int, f func(ctx context.Context) (K, error)) error {
	waitKey := fmt.Sprintf("%s:wait", key)
	cmd := client.SetNX(ctx, waitKey, true, retryT)
	val, err := cmd.Result()
	if err != nil {
		err = fmt.Errorf("goxy:SharedFetch:%s-lock:can't fetch lock: %v", key, err)
		logrus.Error(err)
		return err
	}
	if val == false {
		return fmt.Errorf("goxy:SharedFetch:%s-lock: locked", key)
	}
	var res K
	for i := 0; i < retry; i++ {
		res, err = f(ctx)
		if err == nil {
			break
		}
	}
	if err != nil {
		cmd := client.SetNX(ctx, waitKey, false, retryT)
		if val, lErr := cmd.Result(); err != nil || val == true {
			logrus.Errorf("goxy:SharedFetch:%s-unlock:%s", key, lErr)
		}
		logrus.Errorf("goxy:SharedFetch:fetch(after %d times):%s:%s", retry, key, err.Error())
		return err
	}
	redisSetVal, err := goxy.Encode(res)
	if err != nil {
		logrus.Errorf("goxy:SharedFetch:%s-encode:%s", key, err.Error())
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

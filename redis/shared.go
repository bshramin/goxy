package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/bshramin/goxy"
	"github.com/go-redis/redis/v8"
)

// this is used when multi instances of a code try to use shared cache on redis
// and result fetch function has too much load expensive
// it will return the same result for all of them
// it will checks if the value exists it will return it
// if not it will check if any one else is trying to fetch value
// if it was the case it will not try to fetch
// but if no one of instances has tried to fetch the result it start to fetch and set a waitKey
// too inform others
func SharedFetch[K any](ctx context.Context, client *redis.Client, key string, t time.Duration, f func() (K, error)) (K, error) {
	var res K
	redisRes := client.Get(ctx, key)
	go fetch(ctx, client, key, t, f)
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

func fetch[K any](ctx context.Context, client *redis.Client, key string, t time.Duration, f func() (K, error)) {
	var res K
	waitKey := fmt.Sprintf("%s:wait", key)
	waitTime := t / 2
	waitRes := client.Get(ctx, waitKey)
	_, err := waitRes.Result()
	if waitRes.Err() == redis.Nil && err == nil {
		return
	}
	res, err = f()
	if err != nil {
		return
	}
	redisSetVal, err := goxy.Encode(res)
	if err != nil {
		return
	}
	err = client.Set(ctx, key, redisSetVal, t).Err()
	if err == nil {
		client.Set(ctx, waitKey, "", waitTime)
	}
}

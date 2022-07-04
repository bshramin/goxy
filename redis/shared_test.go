package redis

import (
	"context"
	"testing"
	"time"

	"github.com/bshramin/goxy"
	"github.com/go-redis/redismock/v8"
	"github.com/stretchr/testify/assert"
)

const (
	key      = "testKey"
	waitKey  = "testKey:wait"
	waitData = "waitData"
)

type Template struct {
	Name string
}

func TestSharedFetchCreateMap(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	data := Template{"testData"}
	tt := 2 * time.Second
	f := func(ctx context.Context) (Template, error) {
		return data, nil
	}

	redisDb, _ := redismock.NewClientMock()
	_, err := SharedFetch(ctx, redisDb, key, tt, tt/2, 1, f)
	<-time.After(time.Second)
	assert.Error(t, err)
}

// value : false
// wait : false
func TestSharedFetch1(t *testing.T) {
	t.Parallel()
	redisDb, redisMock := redismock.NewClientMock()

	ctx := context.Background()
	data := Template{"testData"}
	tt := 2 * time.Second
	f := func(ctx context.Context) (Template, error) {
		return data, nil
	}
	expSet, _ := goxy.Encode(data)
	redisMock.MatchExpectationsInOrder(false)
	redisMock.ExpectSet(key, expSet, tt).SetVal("OK")
	redisMock.ExpectSetNX(waitKey, true, tt/2).SetVal(true)
	_, err := SharedFetch(ctx, redisDb, key, tt, tt/2, 1, f)
	assert.Error(t, err)
	<-time.After(time.Second)
	assert.NoError(t, redisMock.ExpectationsWereMet())
}

// value : false
// wait : true
func TestSharedFetch2(t *testing.T) {
	t.Parallel()
	tt := 10 * time.Second
	redisDb, redisMock := redismock.NewClientMock()

	ctx := context.Background()
	data := Template{"testData"}
	f := func(ctx context.Context) (Template, error) {
		return data, nil
	}
	expSet, _ := goxy.Encode(data)
	redisMock.MatchExpectationsInOrder(false)
	redisMock.ExpectSet(key, expSet, tt).SetVal("OK")
	redisMock.ExpectSetNX(waitKey, waitData, tt/2).SetVal(true)
	_, err := SharedFetch(ctx, redisDb, key, tt, tt/2, 1, f)
	assert.Error(t, err)
	<-time.After(time.Second)
	assert.Error(t, redisMock.ExpectationsWereMet())
}

// value : true
// wait : true
func TestSharedFetch3(t *testing.T) {
	t.Parallel()
	tt := 10 * time.Second
	redisDb, redisMock := redismock.NewClientMock()

	ctx := context.Background()
	data := Template{"testData"}
	f := func(ctx context.Context) (Template, error) {
		return data, nil
	}
	expSet, _ := goxy.Encode(data)
	redisMock.MatchExpectationsInOrder(false)
	redisMock.ExpectGet(key).SetVal(expSet)
	redisMock.ExpectSet(key, expSet, tt).SetVal("OK")
	redisMock.ExpectSetNX(waitKey, waitData, tt/2).SetVal(true)
	_, err := SharedFetch(ctx, redisDb, key, tt, tt/2, 1, f)
	assert.NoError(t, err)
	<-time.After(time.Second)
	assert.Error(t, redisMock.ExpectationsWereMet())
}

// value : true
// wait : false
func TestSharedFetch4(t *testing.T) {
	t.Parallel()
	tt := 10 * time.Second
	redisDb, redisMock := redismock.NewClientMock()

	ctx := context.Background()
	data := Template{"testData"}
	f := func(ctx context.Context) (Template, error) {
		return data, nil
	}
	expSet, _ := goxy.Encode(data)
	redisMock.MatchExpectationsInOrder(false)
	redisMock.ExpectSet(key, expSet, tt).SetVal("OK")
	redisMock.ExpectGet(key).SetVal(expSet)
	redisMock.ExpectSetNX(waitKey, true, tt/2).SetVal(true)
	_, err := SharedFetch(ctx, redisDb, key, tt, tt/2, 1, f)
	assert.NoError(t, err)
	<-time.After(time.Second)
	assert.NoError(t, redisMock.ExpectationsWereMet())
}

package redis

import (
	"context"
	"testing"
	"time"

	"github.com/bshramin/goxy"
	"github.com/go-redis/redismock/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
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
	f := func() (Template, error) {
		return data, nil
	}

	redisDb, _ := redismock.NewClientMock()
	_, err := SharedFetch(ctx, redisDb, key, tt, tt/2, f, 1)
	<-time.After(time.Second)
	assert.Error(t, err)
	assert.NotEqual(t, len(keysList), 0)
}

// value : false
// wait : false
func TestSharedFetch1(t *testing.T) {
	t.Parallel()
	redisDb, redisMock := redismock.NewClientMock()
	pool := goredis.NewPool(redisDb)
	rs := redsync.New(pool)
	lock := rs.NewMutex(waitKey, redsync.WithGenValueFunc(func() (string, error) {
		return waitData, nil
	}))
	keysList[redisDb] = make(map[string]*redsync.Mutex)
	keysList[redisDb][key] = lock

	ctx := context.Background()
	data := Template{"testData"}
	tt := 2 * time.Second
	f := func() (Template, error) {
		return data, nil
	}
	expSet, _ := goxy.Encode(data)
	redisMock.MatchExpectationsInOrder(false)
	redisMock.ExpectSet(key, expSet, tt).SetVal("OK")
	redisMock.ExpectSetNX(waitKey, waitData, tt/2).SetVal(true)
	_, err := SharedFetch(ctx, redisDb, key, tt, tt/2, f, 1)
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
	pool := goredis.NewPool(redisDb)
	rs := redsync.New(pool)
	lock := rs.NewMutex(waitKey, redsync.WithGenValueFunc(func() (string, error) {
		return waitData, nil
	}), redsync.WithExpiry(tt/2))
	keysList[redisDb] = make(map[string]*redsync.Mutex)
	keysList[redisDb][key] = lock

	ctx := context.Background()
	data := Template{"testData"}
	f := func() (Template, error) {
		return data, nil
	}
	expSet, _ := goxy.Encode(data)
	redisMock.MatchExpectationsInOrder(false)
	redisMock.ExpectSet(key, expSet, tt).SetVal("OK")
	redisMock.ExpectSetNX(waitKey, waitData, tt/2).SetVal(true)
	err := lock.Lock()
	if err != nil {
		t.Fatal(err)
	}
	_, err = SharedFetch(ctx, redisDb, key, tt, tt/2, f, 1)
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
	pool := goredis.NewPool(redisDb)
	rs := redsync.New(pool)
	lock := rs.NewMutex(waitKey, redsync.WithGenValueFunc(func() (string, error) {
		return waitData, nil
	}), redsync.WithExpiry(tt/2))
	keysList[redisDb] = make(map[string]*redsync.Mutex)
	keysList[redisDb][key] = lock

	ctx := context.Background()
	data := Template{"testData"}
	f := func() (Template, error) {
		return data, nil
	}
	expSet, _ := goxy.Encode(data)
	redisMock.MatchExpectationsInOrder(false)
	redisMock.ExpectGet(key).SetVal(expSet)
	redisMock.ExpectSet(key, expSet, tt).SetVal("OK")
	redisMock.ExpectSetNX(waitKey, waitData, tt/2).SetVal(true)
	err := lock.Lock()
	if err != nil {
		t.Fatal(err)
	}
	_, err = SharedFetch(ctx, redisDb, key, tt, tt/2, f, 1)
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
	pool := goredis.NewPool(redisDb)
	rs := redsync.New(pool)
	lock := rs.NewMutex(waitKey, redsync.WithGenValueFunc(func() (string, error) {
		return waitData, nil
	}), redsync.WithExpiry(tt/2))
	keysList[redisDb] = make(map[string]*redsync.Mutex)
	keysList[redisDb][key] = lock

	ctx := context.Background()
	data := Template{"testData"}
	f := func() (Template, error) {
		return data, nil
	}
	expSet, _ := goxy.Encode(data)
	redisMock.MatchExpectationsInOrder(false)
	redisMock.ExpectSet(key, expSet, tt).SetVal("OK")
	redisMock.ExpectGet(key).SetVal(expSet)
	redisMock.ExpectSetNX(waitKey, waitData, tt/2).SetVal(true)
	_, err := SharedFetch(ctx, redisDb, key, tt, tt/2, f, 1)
	assert.NoError(t, err)
	<-time.After(time.Second)
	assert.NoError(t, redisMock.ExpectationsWereMet())
}

func TestFetchRace(t *testing.T) {
	t.Parallel()
	redisDb, redisMock := redismock.NewClientMock()
	pool := goredis.NewPool(redisDb)
	rs := redsync.New(pool)
	lock := rs.NewMutex(waitKey, redsync.WithGenValueFunc(func() (string, error) {
		return waitData, nil
	}))
	keysList[redisDb] = make(map[string]*redsync.Mutex)
	keysList[redisDb][key] = lock

	ctx := context.Background()
	data := Template{"testData"}
	tt := 2 * time.Second
	f := func() (Template, error) {
		return data, nil
	}
	expSet, _ := goxy.Encode(data)
	redisMock.MatchExpectationsInOrder(false)
	redisMock.ExpectSet(key, expSet, tt).SetVal("OK")
	redisMock.ExpectSetNX(waitKey, waitData, tt/2).SetVal(true)
	var err1, err2 error
	go func() {
		err1 = fetch(ctx, redisDb, key, tt, tt/2, f, 1)
	}()
	<-time.After(time.Second / 10)
	go func() {
		err2 = fetch(ctx, redisDb, key, tt, tt/2, f, 1)
	}()
	<-time.After(time.Second)
	assert.NotEqual(t, err1, err2)
}

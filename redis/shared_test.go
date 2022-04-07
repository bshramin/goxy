package redis

import (
	"context"
	"testing"
	"time"

	"github.com/bshramin/goxy"
	"github.com/go-redis/redismock/v8"
	"github.com/stretchr/testify/assert"
)

// value : false
// wait : false
func TestSharedFetch1(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	key := "testKey"
	waitKey := "testKey:wait"
	type Template struct {
		Name string
	}
	data := Template{"testData"}
	tt := 10 * time.Second
	f := func() (Template, error) {
		return data, nil
	}
	expSet, _ := goxy.Encode(data)

	redisDb, redisMock := redismock.NewClientMock()
	redisMock.MatchExpectationsInOrder(false)
	redisMock.ExpectSet(key, expSet, tt).SetVal("OK")
	redisMock.ExpectSet(waitKey, "", tt/2).SetVal("OK")
	_, err := SharedFetch(ctx, redisDb, key, tt, f)
	assert.Error(t, err)
	<-time.After(time.Second)
	assert.NoError(t, redisMock.ExpectationsWereMet())
}

// value : false
// wait : true
func TestSharedFetch2(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	key := "testKey"
	waitKey := "testKey:wait"
	type Template struct {
		Name string
	}
	data := Template{"testData"}
	tt := 10 * time.Second
	f := func() (Template, error) {
		return data, nil
	}
	// expSet, _ := goxy.Encode(data)

	redisDb, redisMock := redismock.NewClientMock()
	redisMock.MatchExpectationsInOrder(false)
	redisMock.ExpectGet(waitKey).SetVal("")
	_, err := SharedFetch(ctx, redisDb, key, tt, f)
	assert.Error(t, err)
	<-time.After(time.Second)
	assert.NoError(t, redisMock.ExpectationsWereMet())
}

// value : true
// wait : true
func TestSharedFetch3(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	key := "testKey"
	waitKey := "testKey:wait"
	type Template struct {
		Name string
	}
	data := Template{"testData"}
	tt := 10 * time.Second
	f := func() (Template, error) {
		return data, nil
	}
	expSet, _ := goxy.Encode(data)

	redisDb, redisMock := redismock.NewClientMock()
	redisMock.MatchExpectationsInOrder(false)
	redisMock.ExpectGet(key).SetVal(expSet)
	redisMock.ExpectGet(waitKey).SetVal("")
	res, err := SharedFetch(ctx, redisDb, key, tt, f)
	assert.NoError(t, err)
	assert.Equal(t, data, res)
	<-time.After(time.Second)
	assert.NoError(t, redisMock.ExpectationsWereMet())
}

// value : true
// wait : true
func TestSharedFetch4(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	key := "testKey"
	waitKey := "testKey:wait"
	type Template struct {
		Name string
	}
	data := Template{"testData"}
	tt := 10 * time.Second
	f := func() (Template, error) {
		return data, nil
	}
	expSet, _ := goxy.Encode(data)

	redisDb, redisMock := redismock.NewClientMock()
	redisMock.MatchExpectationsInOrder(false)
	redisMock.ExpectGet(key).SetVal(expSet)
	redisMock.ExpectSet(key, expSet, tt).SetVal("OK")
	redisMock.ExpectSet(waitKey, "", tt/2).SetVal("OK")
	res, err := SharedFetch(ctx, redisDb, key, tt, f)
	assert.NoError(t, err)
	assert.Equal(t, data, res)
	<-time.After(time.Second)
	assert.NoError(t, redisMock.ExpectationsWereMet())
}

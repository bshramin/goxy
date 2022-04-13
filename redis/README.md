# Redis

Helpers for Go projects.

### Connection
You can make redis connection and get client to work with.
It works for Single and cluster Redis. 

```golang
// Single connection
client := GetNormalConnection(ctx , host, password , database, timeout, poolSize)

// Cluster connection
client := GetClusterConnection(ctx, host, timeout)
```

### Shared Redis
For a Redis server that is shared with same instances of a code to cache data that has high fetch time or high load, you can use sharedFetch. 

```golang
data, err := SharedFetch(ctx, client, key, duration, retryDuration, dataFetchRetry, dataFetchFunc)
```
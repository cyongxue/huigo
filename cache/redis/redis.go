package redis

import (
    "huigo/cache"
    "github.com/garyburd/redigo/redis"
)

const (
    AdapterNameRedis = "redis"
    DefaultKey = "redis"
)

type CacheRedis struct {
    p           *redis.Pool
    connInfo    string
    dbNum       int
    key         string
    password    string
}

// connect connect to redis
func (r *CacheRedis) connect() {

}



func NewRedisCache() cache.Instance {
    return &CacheRedis{key: DefaultKey}
}

func init()  {
    cache.Register(AdapterNameRedis, NewRedisCache)
}

/**
 * @Time: 2020/3/30 10:38
 * @Author: solacowa@gmail.com
 * @File: client
 * @Software: GoLand
 */

package redis

import (
	"errors"
	"strings"
	"time"

	"github.com/go-redis/redis"
)

type RedisInterface interface {
	Set(k string, v interface{}, expir ...time.Duration) (err error)
	Get(k string) (v string, err error)
	Del(k string) (err error)
	Exists(keys ...string) int64
	HSet(k string, field string, v interface{}) (err error)
	HGet(k string, field string) (res string, err error)
	HGetAll(k string) (res map[string]string, err error)
	HLen(k string) (res int64, err error)
	ZCard(k string) (res int64, err error)
	ZRangeWithScores(k string, start, stop int64) (res []redis.Z, err error)
	ZAdd(k string, score float64, member interface{}) (err error)
	HDelAll(k string) (err error)
	HDel(k string, field string) (err error)
	Keys(pattern string) (res []string, err error)
	Close() error
	Subscribe(channels ...string) *redis.PubSub
	Publish(channel string, message interface{}) error
	Incr(key string, exp time.Duration) error
	SetPrefix(prefix string) RedisInterface
	TTL(key string) time.Duration
}

const (
	RedisCluster = "cluster"
	RedisSingle  = "single"
	expiration   = 600 * time.Second
)

func NewRedisClient(drive, hosts, password, prefix string, db int) (RedisInterface, error) {
	if strings.EqualFold(drive, RedisCluster) {
		return NewRedisCluster(
			strings.Split(hosts, ","),
			password,
			prefix,
		), nil
	} else if strings.EqualFold(drive, RedisSingle) {
		return NewRedisSingle(
			hosts,
			password,
			prefix,
			db,
		), nil
	}

	return nil, errors.New("redis drive is nil!")
}

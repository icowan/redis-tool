/**
 * @Time: 2020/3/30 10:39
 * @Author: solacowa@gmail.com
 * @File: single
 * @Software: GoLand
 */

package redis

import (
	"encoding/json"
	"time"

	"github.com/go-redis/redis"
)

type single struct {
	client *redis.Client
	prefix string
}

func (c *single) LPush(key string, val interface{}) (err error) {
	return c.client.LPush(key, val).Err()
}

func (c *single) RPop(key string) (res string, err error) {
	return c.client.RPop(key).Result()
}

func (c *single) LLen(key string) int64 {
	return c.client.LLen(key).Val()
}

func (c *single) TypeOf(key string) (res string, err error) {
	return c.client.Type(key).Result()
}

func (c *single) Keys(pattern string) (res []string, err error) {
	return c.client.Keys(pattern).Result()
}

func (c *single) ZAdd(k string, score float64, member interface{}) (err error) {
	return c.client.ZAdd(c.setPrefix(k), redis.Z{
		Score:  score,
		Member: member,
	}).Err()
}

func (c *single) ZCard(k string) (res int64, err error) {
	return c.client.ZCard(c.setPrefix(k)).Result()
}

func (c *single) ZRangeWithScores(k string, start, stop int64) (res []redis.Z, err error) {
	return c.client.ZRangeWithScores(c.setPrefix(k), start, stop).Result()
}

func (c *single) HLen(k string) (res int64, err error) {
	return c.client.HLen(c.setPrefix(k)).Result()
}

func (c *single) HGetAll(k string) (res map[string]string, err error) {
	return c.client.HGetAll(c.setPrefix(k)).Result()
}

func (c *single) Exists(keys ...string) int64 {
	return c.client.Exists(keys...).Val()
}

func (c *single) TTL(key string) time.Duration {
	return c.client.TTL(key).Val()
}

func NewRedisSingle(host, password, prefix string, db int) RedisInterface {
	client := redis.NewClient(&redis.Options{
		Addr:     host,
		Password: password, // no password set
		DB:       db,       // use default DB
	})

	return &single{client: client, prefix: prefix}
}

func (c *single) Set(k string, v interface{}, expir ...time.Duration) (err error) {
	var val string
	switch v.(type) {
	case string:
		val = v.(string)
	default:
		b, _ := json.Marshal(v)
		val = string(b)
	}

	exp := expiration
	if len(expir) == 1 {
		exp = expir[0]
	}

	return c.client.Set(c.setPrefix(k), val, exp).Err()
}

func (c *single) Get(k string) (v string, err error) {
	return c.client.Get(c.setPrefix(k)).Result()
}

func (c *single) Del(k string) (err error) {
	return c.client.Del(c.setPrefix(k)).Err()
}

func (c *single) HSet(k string, field string, v interface{}) (err error) {
	var val string
	switch v.(type) {
	case string:
		val = v.(string)
	default:
		b, _ := json.Marshal(v)
		val = string(b)
	}
	return c.client.HSet(c.setPrefix(k), field, val).Err()
}

func (c *single) HGet(k string, field string) (res string, err error) {
	return c.client.HGet(c.setPrefix(k), field).Result()
}

func (c *single) HDelAll(k string) (err error) {
	res, err := c.client.HKeys(c.setPrefix(k)).Result()
	if err != nil {
		return
	}
	return c.client.HDel(c.setPrefix(k), res...).Err()
}

func (c *single) HDel(k string, field string) (err error) {
	return c.client.HDel(c.setPrefix(k), field).Err()
}

func (c *single) setPrefix(s string) string {
	return c.prefix + s
}

func (c *single) Close() error {
	return c.client.Close()
}

func (c *single) Subscribe(channels ...string) *redis.PubSub {
	return c.client.Subscribe(channels...)
}

func (c *single) Publish(channel string, message interface{}) error {
	return c.client.Publish(channel, message).Err()
}

func (c *single) Incr(key string, expiration time.Duration) error {
	defer func() {
		c.client.Expire(c.setPrefix(key), expiration)
	}()
	return c.client.Incr(c.setPrefix(key)).Err()
}

func (c *single) SetPrefix(prefix string) RedisInterface {
	c.prefix = prefix
	return c
}

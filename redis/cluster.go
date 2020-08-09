/**
 * @Time: 2020/3/30 10:39
 * @Author: solacowa@gmail.com
 * @File: cluster
 * @Software: GoLand
 */

package redis

import (
	"encoding/json"
	"time"

	"github.com/go-redis/redis"
)

type cluster struct {
	client *redis.ClusterClient
	//prefix func(s string) string
	prefix string
}

func (c *cluster) LPush(key string, val interface{}) (err error) {
	return c.client.LPush(key, val).Err()
}

func (c *cluster) RPop(key string) (res string, err error) {
	return c.client.RPop(key).Result()
}

func (c *cluster) LLen(key string) int64 {
	return c.client.LLen(key).Val()
}

func (c *cluster) TypeOf(key string) (res string, err error) {
	return c.client.Type(key).Result()
}

func (c *cluster) Keys(pattern string) (res []string, err error) {
	return c.client.Keys(pattern).Result()
}

func (c *cluster) ZAdd(k string, score float64, member interface{}) (err error) {
	return c.client.ZAdd(c.setPrefix(k), redis.Z{
		Score:  score,
		Member: member,
	}).Err()
}

func (c *cluster) ZRangeWithScores(k string, start, stop int64) (res []redis.Z, err error) {
	return c.client.ZRangeWithScores(c.setPrefix(k), start, stop).Result()
}

func (c *cluster) ZCard(k string) (res int64, err error) {
	return c.client.ZCard(c.setPrefix(k)).Result()
}

func (c *cluster) HLen(k string) (res int64, err error) {
	return c.client.HLen(c.setPrefix(k)).Result()
}

func (c *cluster) HGetAll(k string) (res map[string]string, err error) {
	return c.client.HGetAll(c.setPrefix(k)).Result()
}

func (c *cluster) Incr(key string, exp time.Duration) error {
	defer func() {
		c.client.Expire(c.setPrefix(key), exp)
	}()
	return c.client.Incr(c.setPrefix(key)).Err()
}

func NewRedisCluster(hosts []string, password, prefix string) RedisInterface {
	return &cluster{client: redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    hosts,
		Password: password,
	}), prefix: prefix + ""}
}

func (c *cluster) Exists(keys ...string) int64 {
	return c.client.Exists(keys...).Val()
}

func (c *cluster) TTL(key string) time.Duration {
	return c.client.TTL(key).Val()
}

func (c *cluster) Set(k string, v interface{}, expir ...time.Duration) (err error) {
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

func (c *cluster) Get(k string) (v string, err error) {
	return c.client.Get(c.setPrefix(k)).Result()
}

func (c *cluster) Del(k string) (err error) {
	return c.client.Del(c.setPrefix(k)).Err()
}

func (c *cluster) HSet(k string, field string, v interface{}) (err error) {
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

func (c *cluster) HGet(k string, field string) (res string, err error) {
	return c.client.HGet(c.setPrefix(k), field).Result()
}

func (c *cluster) HDelAll(k string) (err error) {
	res, err := c.client.HKeys(c.setPrefix(k)).Result()
	if err != nil {
		return
	}
	return c.client.HDel(c.setPrefix(k), res...).Err()
}

func (c *cluster) HDel(k string, field string) (err error) {
	return c.client.HDel(c.setPrefix(k), field).Err()
}

func (c *cluster) setPrefix(s string) string {
	return c.prefix + s
}

func (c *cluster) Close() error {
	return c.client.Close()
}

func (c *cluster) Subscribe(channels ...string) *redis.PubSub {
	return c.client.Subscribe(channels...)
}

func (c *cluster) Publish(channel string, message interface{}) error {
	return c.client.Publish(channel, message).Err()
}

func (c *cluster) SetPrefix(prefix string) RedisInterface {
	c.prefix = prefix
	return c
}

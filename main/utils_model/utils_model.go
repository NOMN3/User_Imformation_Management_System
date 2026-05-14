package utils

import (
	"BOOK/main/utils_model/redis"
	"time"
)

// 以下是存普通key—values,key是文章_id，value是文章内容

func Set(key, value string, expiration time.Duration) error {
	return redis.RedisClient.Set(redis.Ctx, key, value, expiration).Err()
}

func Get(key string) (string, error) {
	return redis.RedisClient.Get(redis.Ctx, key).Result()
}

func Delete(key string) error {
	return redis.RedisClient.Del(redis.Ctx, key).Err()
}

// 以下是存list，key是：用户_id，values是文章_id

// 存数据
func ListPush(key string, values ...string) error {
	return redis.RedisClient.RPush(redis.Ctx, key, values).Err()
}

// 删除
func ListRemove(key string, values ...string) error {
	return redis.RedisClient.LRem(redis.Ctx, key, 0, values).Err()
}

// 有范围地查询
func ListRange(key string, start, stop int64) ([]string, error) {
	return redis.RedisClient.LRange(redis.Ctx, key, start, stop).Result()
}

// 查有多少个
func ListSize(key string) (int64, error) {
	return redis.RedisClient.LLen(redis.Ctx, key).Result()
}

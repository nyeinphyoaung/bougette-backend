package utilities

import (
	"bougette-backend/configs"
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisClient() *RedisClient {
	redis := &RedisClient{
		client: configs.Envs.Redis,
		ctx:    context.Background(),
	}
	return redis
}

func (r *RedisClient) Set(key string, value interface{}, expiration time.Duration) error {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return r.client.Set(r.ctx, key, jsonValue, expiration).Err()
}

func (r *RedisClient) Get(key string) (string, error) {
	return r.client.Get(r.ctx, key).Result()
}

func (r *RedisClient) GetObject(key string, dest interface{}) error {
	val, err := r.client.Get(r.ctx, key).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(val), dest)
}

func (r *RedisClient) Delete(key string) error {
	return r.client.Del(r.ctx, key).Err()
}

func (r *RedisClient) Exists(key string) (bool, error) {
	count, err := r.client.Exists(r.ctx, key).Result()
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *RedisClient) SetString(key string, value string, expiration time.Duration) error {
	return r.client.Set(r.ctx, key, value, expiration).Err()
}

func (r *RedisClient) GetString(key string) (string, error) {
	return r.client.Get(r.ctx, key).Result()
}

func (r *RedisClient) Increment(key string) (int64, error) {
	return r.client.Incr(r.ctx, key).Result()
}

func (r *RedisClient) Decrement(key string) (int64, error) {
	return r.client.Decr(r.ctx, key).Result()
}

func (r *RedisClient) SetHash(key, field string, value interface{}) error {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return r.client.HSet(r.ctx, key, field, jsonValue).Err()
}

func (r *RedisClient) GetHash(key, field string) (string, error) {
	return r.client.HGet(r.ctx, key, field).Result()
}

func (r *RedisClient) GetHashAll(key string) (map[string]string, error) {
	return r.client.HGetAll(r.ctx, key).Result()
}

func (r *RedisClient) AddToSet(key, member string) error {
	return r.client.SAdd(r.ctx, key, member).Err()
}

func (r *RedisClient) RemoveFromSet(key, member string) error {
	return r.client.SRem(r.ctx, key, member).Err()
}

func (r *RedisClient) IsMemberOfSet(key, member string) (bool, error) {
	return r.client.SIsMember(r.ctx, key, member).Result()
}

func (r *RedisClient) PushToList(key, value string) error {
	return r.client.LPush(r.ctx, key, value).Err()
}

func (r *RedisClient) PopFromList(key string) (string, error) {
	return r.client.RPop(r.ctx, key).Result()
}

func (r *RedisClient) GetListRange(key string, start, stop int64) ([]string, error) {
	return r.client.LRange(r.ctx, key, start, stop).Result()
}

func (r *RedisClient) Client() *redis.Client {
	return r.client
}

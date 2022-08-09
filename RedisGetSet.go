package Tool

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/vmihailenco/msgpack/v5"
)

func RedisGetBytes(c context.Context, rds *redis.Client, key string) (bytes []byte, err error) {
	cmd := rds.Get(c, key)
	return cmd.Bytes()
}
func RedisSetBytes(c context.Context, rds *redis.Client, key string, param []byte, expiration time.Duration) (err error) {
	cmd := rds.Set(c, key, param, expiration)
	if cmd.Err() != nil {
		return cmd.Err()
	}
	return nil
}
func RedisGetString(c context.Context, rds *redis.Client, key string) (string, error) {
	cmd := rds.Get(c, key)
	return cmd.String(), cmd.Err()
}
func RedisSetString(c context.Context, rds *redis.Client, key string, param string, expiration time.Duration) (err error) {
	cmd := rds.Set(c, key, param, expiration)
	if cmd.Err() != nil {
		return cmd.Err()
	}
	return nil
}
func RedisGet(c context.Context, rds *redis.Client, key string, param interface{}) (err error) {
	cmd := rds.Get(c, key)
	data, err := cmd.Bytes()
	if err != nil {
		return err
	}
	return msgpack.Unmarshal(data, param)
}
func RedisSet(c context.Context, rds *redis.Client, key string, param interface{}, expiration time.Duration) (err error) {
	bytes, err := msgpack.Marshal(param)
	if err != nil {
		return err
	}
	status := rds.Set(c, key, bytes, expiration)
	return status.Err()
}
func RedisHGet(c context.Context, rds *redis.Client, key string, field string, param interface{}) (err error) {
	cmd := rds.HGet(c, key, field)
	data, err := cmd.Bytes()
	if err != nil {
		return err
	}
	return msgpack.Unmarshal(data, param)
}

func RedisHSet(c context.Context, rds *redis.Client, key string, field string, param interface{}) (err error) {
	bytes, err := msgpack.Marshal(param)
	if err != nil {
		return err
	}
	status := rds.HSet(c, key, field, bytes)
	return status.Err()
}

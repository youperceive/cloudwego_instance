package redis

import (
	"context"
	"log"
	"os"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

var rdb *redis.Client

func init() {
	rdb = redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Println("failed to connect to redis")
	}
}

func MakeKey(params []string) string {
	return strings.Join(params, ":")
}

func SetWithCount(ctx context.Context, key, code string, expire time.Duration, maxCount int) error {
	err := rdb.HMSet(ctx, key, map[string]any{
		"code":   code,
		"remain": maxCount,
	}).Err()
	if err != nil {
		return err
	}
	return rdb.Expire(ctx, key, expire).Err()
}

func GetAndDecrementCount(ctx context.Context, key string) (string, int, error) {
	remain, err := rdb.HGet(ctx, key, "remain").Int()
	if err != nil {
		return "", 0, err
	}
	if remain <= 0 {
		return "", 0, redis.Nil
	}

	code, err := rdb.HGet(ctx, key, "code").Result()
	if err != nil {
		return "", 0, err
	}

	newRemain := remain - 1
	if newRemain > 0 {
		err := rdb.HSet(ctx, key, "remain", newRemain).Err()
		if err != nil {
			return "", 0, err
		}
	} else {
		err := Delete(ctx, key)
		if err != nil {
			return "", 0, err
		}
	}

	return code, newRemain, nil
}

func Delete(ctx context.Context, key string) error {
	return rdb.Del(ctx, key).Err()
}

func Exists(ctx context.Context, key string) (bool, error) {
	count, err := rdb.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

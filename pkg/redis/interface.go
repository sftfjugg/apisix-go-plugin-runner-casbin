package redis

import (
	"github.com/go-redis/redis"
	"time"
)

type Interface interface {
	// Keys retrieves all keys match the given pattern
	Keys(pattern string) ([]string, error)

	// Get retrieves the value of the given key, return error if key doesn't exist
	Get(key string) (string, error)

	// Set sets the value and living duration of the given key, zero duration means never expire
	Set(key string, value string, duration time.Duration) error

	// Del deletes the given key, no error returned if the key doesn't exists
	Del(keys ...string) error

	// LPush Insert one or more values into the header of the list key. If there are multiple values, each value is inserted into the header from left to right.
	LPush(key string, values ...interface{}) error

	// RPop Remove and return the last element of the list key.
	RPop(key string) (string, error)

	// BRPop Remove and return the last element of the list key.
	BRPop(keys ...string) ([]string, error)

	Publish(key string, message interface{}) (int64, error)

	Subscribe(keys ...string) *redis.PubSub

	// Exists checks the existence of a give key
	Exists(keys ...string) (bool, error)

	// Expire updates object's expiration time, return err if key doesn't exist
	Expire(key string, duration time.Duration) error
}

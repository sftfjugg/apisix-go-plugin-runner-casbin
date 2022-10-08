package redis

import (
	"fmt"
	"github.com/apache/apisix-go-plugin-runner/pkg/log"
	"github.com/go-redis/redis"
	"github.com/spf13/viper"
	"time"
)

var rc Interface
var stopChRedis chan struct{}

type Client struct {
	client *redis.Client
}

func Setup() {
	var err error
	stopChRedis = make(chan struct{})
	rc, err = newRedisClient(
		viper.GetString("redis.host"),
		viper.GetInt("redis.port"),
		viper.GetString("redis.password"),
		viper.GetInt("redis.db"),
		stopChRedis,
	)

	if err != nil {
		log.Fatalf("Redis 初始化失败")
	}

	return
}

func newRedisClient(host string, port int, password string, db int, stopCh <-chan struct{}) (Interface, error) {
	var r Client

	redisOptions := &redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: password,
		DB:       db,

		PoolSize:     20,
		MinIdleConns: 10,
	}

	if stopCh == nil {
		log.Fatalf("no stop channel passed, redis connections will leak.")
	}

	r.client = redis.NewClient(redisOptions)

	if err := r.client.Ping().Err(); err != nil {
		_ = r.client.Close()
		return nil, err
	}

	// close redis in case of connection leak
	if stopCh != nil {
		go func() {
			<-stopCh
			if err := r.client.Close(); err != nil {
				log.Errorf(err.Error())
			} else {
				log.Infof("Redis 链接已关闭")
			}
		}()
	}

	return &r, nil
}

func (r *Client) Get(key string) (string, error) {
	return r.client.Get(key).Result()
}

func (r *Client) Keys(pattern string) ([]string, error) {
	return r.client.Keys(pattern).Result()
}

func (r *Client) Set(key string, value string, duration time.Duration) error {
	return r.client.Set(key, value, duration).Err()
}

func (r *Client) Del(keys ...string) error {
	return r.client.Del(keys...).Err()
}

func (r *Client) LPush(key string, values ...interface{}) error {
	return r.client.LPush(key, values...).Err()
}

func (r *Client) RPop(key string) (string, error) {
	return r.client.RPop(key).Result()
}

func (r *Client) BRPop(keys ...string) ([]string, error) {
	return r.client.BRPop(time.Duration(0), keys...).Result()
}

func (r *Client) Exists(keys ...string) (bool, error) {
	existedKeys, err := r.client.Exists(keys...).Result()
	if err != nil {
		return false, err
	}

	return len(keys) == int(existedKeys), nil
}

func (r *Client) Publish(key string, message interface{}) (int64, error) {
	return r.client.Publish(key, message).Result()
}

func (r *Client) Subscribe(keys ...string) *redis.PubSub {
	return r.client.Subscribe(keys...)
}

func (r *Client) Expire(key string, duration time.Duration) error {
	return r.client.Expire(key, duration).Err()
}

func Rc() Interface {
	return rc
}

func StopChRedis() {
	stopChRedis <- struct{}{}
}

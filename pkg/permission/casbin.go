package permission

import (
	"fmt"
	"github.com/apache/apisix-go-plugin-runner/pkg/db"
	"github.com/apache/apisix-go-plugin-runner/pkg/log"
	"github.com/apache/apisix-go-plugin-runner/pkg/redis"
	redis2 "github.com/go-redis/redis"
	"time"

	"github.com/casbin/casbin/v2"
	gormAdapter "github.com/casbin/gorm-adapter/v3"
	"github.com/spf13/viper"
)

const (
	loadKey = "load:casbin:data"
)

var enforcer *casbin.SyncedEnforcer

func Setup() {
	setEnforcer()        // 创建权限实例
	loadPermission()     // 启动插件时首次加载权限
	authLoadPermission() // 定期及watch变化同步权限
}

func loadPermission() {
	err := Enforcer().LoadPolicy()
	if err != nil {
		log.Fatalf("从数据库加载策略失败，错误：%v", err)
	}
}

func authLoadPermission() {
	// 定时同步策略
	if viper.GetBool("casbin.isTiming") {
		// 间隔多长时间同步一次权限策略，单位：秒
		Enforcer().StartAutoLoadPolicy(time.Second * time.Duration(viper.GetInt("casbin.intervalTime")))
	}

	// Watch 权限
	go func() {
		pubsub := redis.Rc().Subscribe(loadKey)
		defer func(pubsub *redis2.PubSub) {
			err := pubsub.Close()
			if err != nil {
				log.Fatalf(err.Error())
			}
		}(pubsub)
		for _ = range pubsub.Channel() {
			loadPermission()
		}
	}()
}

func setEnforcer() {
	var (
		err     error
		adapter *gormAdapter.Adapter
	)
	adapter, err = gormAdapter.NewAdapterByDBWithCustomTable(db.Orm(), nil, viper.GetString("casbin.tableName"))
	if err != nil {
		log.Fatalf("创建 casbin gorm adapter 失败，错误：%v", err)
	}

	enforcer, err = casbin.NewSyncedEnforcer(viper.GetString("casbin.rbacModel"), adapter)
	if err != nil {
		log.Fatalf("创建 casbin enforcer 失败，错误：%v", err)
	}
}

func CheckPermission(obj, act, sub string, isAdmin bool) (ok bool, err error) {
	if isAdmin {
		ok = true
	} else {
		//判断策略中是否存在
		ok, err = Enforcer().Enforce(sub, obj, act)
		if !ok {
			err = fmt.Errorf("the interface cannot be called via %s for the time being %s", act, obj)
			return
		}
	}
	return
}

func Enforcer() *casbin.SyncedEnforcer {
	return enforcer
}

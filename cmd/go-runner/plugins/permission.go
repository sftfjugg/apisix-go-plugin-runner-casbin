package plugins

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/apache/apisix-go-plugin-runner/pkg/jwtauth"
	"github.com/apache/apisix-go-plugin-runner/pkg/permission"
	"net/http"
	"strings"

	pkgHTTP "github.com/apache/apisix-go-plugin-runner/pkg/http"
	"github.com/apache/apisix-go-plugin-runner/pkg/log"
	"github.com/apache/apisix-go-plugin-runner/pkg/plugin"
)

func init() {
	err := plugin.RegisterPlugin(&Permission{})
	if err != nil {
		log.Fatalf("failed to register plugin permission: %s", err)
	}
}

// Permission .
type Permission struct {
	plugin.DefaultPlugin
}

type PermissionConf struct{}

func (p *Permission) Name() string {
	return "permission"
}

func (p *Permission) ParseConf(in []byte) (interface{}, error) {
	conf := PermissionConf{}
	err := json.Unmarshal(in, &conf)
	return conf, err
}

func (p *Permission) parseToken(r pkgHTTP.Request) (claims *jwtauth.Claims, err error) {
	token := r.Header().Get("Authorization")
	if token == "" {
		err = errors.New("not logged in yet")
		return
	}

	// 按空格分割
	parts := strings.SplitN(token, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Bearer") {
		err = errors.New("the token format is incorrect")
		return
	}

	// parts[1]是获取到的tokenString，我们使用之前定义好的解析JWT的函数来解析它
	claims, err = jwtauth.ParseToken(parts[1])
	if err != nil {
		return
	}
	return
}

func (p *Permission) response(code int, message string) (resp []byte) {
	resp, _ = json.Marshal(map[string]interface{}{
		"code":    code,
		"message": message,
	})
	return
}

func (p *Permission) RequestFilter(conf interface{}, w http.ResponseWriter, r pkgHTTP.Request) {

	var (
		claims     *jwtauth.Claims
		err        error
		ok         bool
		statusCode int
	)

	w.Header().Add("X-Gateway", "true")

	// 解析 token
	claims, err = p.parseToken(r)
	if err != nil {
		statusCode = 44000 // 登录异常
		goto write
	}

	// 验证是否有权限
	ok, err = permission.CheckPermission(string(r.Path()), r.Method(), claims.Username, claims.IsAdmin)
	if !ok || err != nil {
		statusCode = 43000 // 无权限
		goto write
	}

write:
	if err != nil {
		resp := p.response(statusCode, fmt.Sprintf("Authentication failed, %s", err.Error()))
		_, err = w.Write(resp)
		if err != nil {
			log.Errorf("failed to write: %s", err)
		}
	}
}

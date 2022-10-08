module github.com/apache/apisix-go-plugin-runner

go 1.15

require (
	github.com/ReneKroon/ttlcache/v2 v2.4.0
	github.com/api7/ext-plugin-proto v0.6.0
	github.com/casbin/casbin/v2 v2.55.1
	github.com/casbin/gorm-adapter/v3 v3.10.0
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/golang-jwt/jwt v3.2.2+incompatible
	github.com/google/flatbuffers v2.0.0+incompatible
	github.com/jackc/pgx/v4 v4.17.2
	github.com/kr/pretty v0.3.0 // indirect
	github.com/lib/pq v1.10.7
	github.com/rogpeppe/go-internal v1.8.0 // indirect
	github.com/spf13/cobra v1.2.1
	github.com/spf13/viper v1.8.1
	github.com/stretchr/testify v1.8.0
	github.com/thediveo/enumflag v0.10.1
	go.uber.org/multierr v1.7.0 // indirect
	go.uber.org/zap v1.17.0
	golang.org/x/tools v0.1.9 // indirect
	google.golang.org/protobuf v1.28.0 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gorm.io/driver/mysql v1.3.6
	gorm.io/driver/postgres v1.3.10
	gorm.io/gorm v1.23.10
)

replace (
	github.com/miekg/dns v1.0.14 => github.com/miekg/dns v1.1.25
	// github.com/thediveo/enumflag@v0.10.1 depends on github.com/spf13/cobra@v0.0.7
	github.com/spf13/cobra v0.0.7 => github.com/spf13/cobra v1.2.1
)

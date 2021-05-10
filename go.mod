module Book

go 1.14

replace github.com/coreos/bbolt v1.3.4 => go.etcd.io/bbolt v1.3.4

replace google.golang.org/grpc v1.29.1 => google.golang.org/grpc v1.26.0

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gin-gonic/gin v1.6.3
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/jinzhu/gorm v1.9.16
	github.com/onsi/gomega v1.12.0 // indirect
	github.com/sirupsen/logrus v1.4.2
	go.etcd.io/etcd/client/v3 v3.5.0-alpha.0
	google.golang.org/genproto v0.0.0-20210312152112-fc591d9ea70f // indirect
	gopkg.in/alexcesaro/quotedprintable.v3 v3.0.0-20150716171945-2caba252f4dc // indirect
	gopkg.in/gomail.v2 v2.0.0-20160411212932-81ebce5c23df
	gopkg.in/yaml.v2 v2.4.0
)

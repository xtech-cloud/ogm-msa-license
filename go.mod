module ogm-msa-license

go 1.16

require (
	github.com/asim/go-micro/plugins/config/encoder/yaml/v3 v3.7.0
	github.com/asim/go-micro/plugins/config/source/etcd/v3 v3.7.0
	github.com/asim/go-micro/plugins/logger/logrus/v3 v3.7.0
	github.com/asim/go-micro/plugins/registry/etcd/v3 v3.7.0
	github.com/asim/go-micro/plugins/server/grpc/v3 v3.7.0
	github.com/asim/go-micro/v3 v3.7.0
	github.com/satori/go.uuid v1.2.0
	github.com/sirupsen/logrus v1.8.1
	github.com/xtech-cloud/omo-msp-license v3.0.0+incompatible
	gopkg.in/yaml.v2 v2.4.0
	gorm.io/driver/mysql v1.1.1
	gorm.io/driver/sqlite v1.1.4
	gorm.io/gorm v1.21.12
)

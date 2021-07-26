module ogm-msa-license

go 1.16

require (
	github.com/asim/go-micro/plugins/config/encoder/yaml/v3 v3.0.0-20210721080634-e1bc7e302871
	github.com/asim/go-micro/plugins/config/source/etcd/v3 v3.0.0-20210721080634-e1bc7e302871
	github.com/asim/go-micro/plugins/logger/logrus/v3 v3.0.0-20210721080634-e1bc7e302871
	github.com/asim/go-micro/plugins/registry/etcd/v3 v3.0.0-20210721080634-e1bc7e302871
	github.com/asim/go-micro/v3 v3.5.2
	github.com/satori/go.uuid v1.2.0
	github.com/sirupsen/logrus v1.8.1
	github.com/xtech-cloud/omo-msp-license v3.0.0+incompatible
	golang.org/x/net v0.0.0-20210614182718-04defd469f4e // indirect
	google.golang.org/genproto v0.0.0-20210721163202-f1cecdd8b78a // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
	gorm.io/driver/mysql v1.1.1
	gorm.io/driver/sqlite v1.1.4
	gorm.io/gorm v1.21.12
)

APP_NAME := ogm-license
BUILD_VERSION   := $(shell git tag --contains)
BUILD_TIME      := $(shell date "+%F %T")
COMMIT_SHA1     := $(shell git rev-parse HEAD )

SPACE_KEY := $(shell cat /tmp/msa-license-key)
SPACE_SECRET := $(shell cat /tmp/msa-license-secret)
NUMBER := $(shell cat /tmp/msa-license-number)

.PHONY: build
build:
	go build -ldflags \
		"\
		-X 'main.BuildVersion=${BUILD_VERSION}' \
		-X 'main.BuildTime=${BUILD_TIME}' \
		-X 'main.CommitID=${COMMIT_SHA1}' \
		"\
		-o ./bin/${APP_NAME}

.PHONY: run
run:
	./bin/${APP_NAME}

.PHONY: install
install:
	go install

.PHONY: clean
clean:
	rm -rf /tmp/msa-license.db

.PHONY: call
call:
	gomu --registry=etcd --client=grpc call xtc.ogm.license Healthy.Echo '{"msg":"hello"}'
	# 创建空间
	gomu --registry=etcd --client=grpc call xtc.ogm.license Space.Create '{"name":"test"}'
	# 查询存在的空间
	gomu --registry=etcd --client=grpc call xtc.ogm.license Space.Query '{"name":"test"}'
	# 查询不存在的空间
	gomu --registry=etcd --client=grpc call xtc.ogm.license Space.Query '{"name":"test-1"}'
	# 列举空间
	gomu --registry=etcd --client=grpc call xtc.ogm.license Space.List '{"offset":0, "count":100}'
	# 错误的参数
	gomu --registry=etcd --client=grpc call xtc.ogm.license Key.Generate '{"space":"1212"}'
	# 使用默认参数生成
	gomu --registry=etcd --client=grpc call xtc.ogm.license Key.Generate '{"space":"test"}'
	# 使用指定参数生成
	gomu --registry=etcd --client=grpc call xtc.ogm.license Key.Generate '{"space":"test", "count":3, "capacity":2,"expiry":4,"storage":"mydata","profile":"myprofile"}'
	# 查询存在的激活码
	gomu --registry=etcd --client=grpc call xtc.ogm.license Key.Query '{"number":"${NUMBER}"}'
	# 缺少参数
	gomu --registry=etcd --client=grpc call xtc.ogm.license Key.Activate '{"number":"${NUMBER}", "consumer":"223344"}'
	# 错误的参数
	gomu --registry=etcd --client=grpc call xtc.ogm.license Key.Activate '{"number":"${NUMBER}", "consumer":"223344", "space":"test2"}'
	# 激活
	gomu --registry=etcd --client=grpc call xtc.ogm.license Key.Activate '{"number":"${NUMBER}", "consumer":"223344", "space":"test"}'
	# 激活已激活过的
	gomu --registry=etcd --client=grpc call xtc.ogm.license Key.Activate '{"number":"${NUMBER}", "consumer":"223344", "space":"test"}'
	# 激活超过capacity的
	gomu --registry=etcd --client=grpc call xtc.ogm.license Key.Activate '{"number":"${NUMBER}", "consumer":"22334455", "space":"test"}'
	# 激活超过capacity的
	gomu --registry=etcd --client=grpc call xtc.ogm.license Key.Activate '{"number":"${NUMBER}", "consumer":"44556677", "space":"test"}'
	# 挂起证书
	gomu --registry=etcd --client=grpc call xtc.ogm.license Key.Suspend '{"space":"test", "number":"${NUMBER}", "ban":1, "reason":"unknown"}'
	# 使用挂起的证书激活
	gomu --registry=etcd --client=grpc call xtc.ogm.license Key.Activate '{"number":"${NUMBER}", "consumer":"223344", "space":"test"}'
	# 恢复挂起的证书
	gomu --registry=etcd --client=grpc call xtc.ogm.license Key.Suspend '{"space":"test", "number":"${NUMBER}", "ban":0, "reason":"unknown"}'
	# 列举激活码
	gomu --registry=etcd --client=grpc call xtc.ogm.license Key.List '{"offset":0, "count":100, "space": "test"}'
	# 拉取证书
	gomu --registry=etcd --client=grpc call xtc.ogm.license Certificate.Pull '{"space":"test", "consumer":"223344"}'

.PHONY: tcall
tcall:
	mkdir -p ./bin
	go build -o ./bin/ ./client
	./bin/client

.PHONY: post
post:
	curl -X POST -H 'Content-Type:application/json' -d '{"msg":"hello"}' 127.0.0.1/ogm/license/Healthy/Echo


.PHONY: dist
dist:
	mkdir dist
	tar -zcf dist/${APP_NAME}-${BUILD_VERSION}.tar.gz ./bin/${APP_NAME}

.PHONY: docker
docker:
	docker build -t xtechcloud/${APP_NAME}:${BUILD_VERSION} .
	docker rm -f ${APP_NAME}
	docker run --restart=always --name=${APP_NAME} --net=host -e MSA_REGISTRY_ADDRESS='localhost:2379' -d xtechcloud/${APP_NAME}:${BUILD_VERSION}
	docker logs -f ${APP_NAME}

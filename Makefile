APP_NAME := omo-msa-license
BUILD_VERSION   := $(shell git tag --contains)
BUILD_TIME      := $(shell date "+%F %T")
COMMIT_SHA1     := $(shell git rev-parse HEAD )


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
SPACE_KEY := $(shell cat /tmp/msa-license-key)
SPACE_SECRET := $(shell cat /tmp/msa-license-secret)
NUMBER := $(shell cat /tmp/msa-license-number)
call:
	# 创建空间
	MICRO_REGISTRY=consul micro call omo.msa.license Space.Create '{"name":"test"}'
	# 查询存在的空间
	MICRO_REGISTRY=consul micro call omo.msa.license Space.Query '{"name":"test"}'
	# 查询不存在的空间
	MICRO_REGISTRY=consul micro call omo.msa.license Space.Query '{"name":"test-1"}'
	# 列举空间
	MICRO_REGISTRY=consul micro call omo.msa.license Space.List '{"offset":0, "count":100}'
	# 错误的参数
	MICRO_REGISTRY=consul micro call omo.msa.license Key.Generate '{"space":"1212"}'
	# 使用默认参数生成
	MICRO_REGISTRY=consul micro call omo.msa.license Key.Generate '{"space":"test"}'
	# 使用指定参数生成
	MICRO_REGISTRY=consul micro call omo.msa.license Key.Generate '{"space":"test", "count":3, "capacity":2,"expiry":4,"storage":"mydata","profile":"myprofile"}'
	# 查询存在的激活码
	MICRO_REGISTRY=consul micro call omo.msa.license Key.Query '{"number":"${NUMBER}"}'
	# 缺少参数
	MICRO_REGISTRY=consul micro call omo.msa.license Key.Activate '{"number":"${NUMBER}", "consumer":"223344"}'
	# 错误的参数
	MICRO_REGISTRY=consul micro call omo.msa.license Key.Activate '{"number":"${NUMBER}", "consumer":"223344", "space":"test2"}'
	# 激活
	MICRO_REGISTRY=consul micro call omo.msa.license Key.Activate '{"number":"${NUMBER}", "consumer":"223344", "space":"test"}'
	# 激活已激活过的
	MICRO_REGISTRY=consul micro call omo.msa.license Key.Activate '{"number":"${NUMBER}", "consumer":"223344", "space":"test"}'
	# 激活超过capacity的
	MICRO_REGISTRY=consul micro call omo.msa.license Key.Activate '{"number":"${NUMBER}", "consumer":"22334455", "space":"test"}'
	# 激活超过capacity的
	MICRO_REGISTRY=consul micro call omo.msa.license Key.Activate '{"number":"${NUMBER}", "consumer":"44556677", "space":"test"}'
	# 挂起证书
	MICRO_REGISTRY=consul micro call omo.msa.license Key.Suspend '{"space":"test", "number":"${NUMBER}", "ban":1, "reason":"unknown"}'
	# 使用挂起的证书激活
	MICRO_REGISTRY=consul micro call omo.msa.license Key.Activate '{"number":"${NUMBER}", "consumer":"223344", "space":"test"}'
	# 恢复挂起的证书
	MICRO_REGISTRY=consul micro call omo.msa.license Key.Suspend '{"space":"test", "number":"${NUMBER}", "ban":0, "reason":"unknown"}'
	# 列举激活码
	MICRO_REGISTRY=consul micro call omo.msa.license Key.List '{"offset":0, "count":100, "space": "test"}'
	# 拉取证书
	MICRO_REGISTRY=consul micro call omo.msa.license Certificate.Pull '{"space":"test", "consumer":"223344"}'

.PHONY: tcall
tcall:
	mkdir -p ./bin
	go build -o ./bin/ ./client
	./bin/client

.PHONY: dist
dist:
	mkdir dist
	tar -zcf dist/${APP_NAME}-${BUILD_VERSION}.tar.gz ./bin/${APP_NAME}

.PHONY: docker
docker:
	docker build . -t omo-msa-startkit:latest

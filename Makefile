BUILD_DIR=./out/bin
BINARY_NAME=cryptomath-captcha
SEED_BINARY_NAME=${BINARY_NAME}-seed
DOCKER_COMPOSE_FILE=docker-compose.yaml

#	Go section
.PHONY: build
build:
	mkdir -p ${BUILD_DIR}
	cd ${BUILD_DIR}
	go build -mod vendor -o ${BUILD_DIR}/${BINARY_NAME} ./cmd/captcha/main.go

.PHONY: build-seed
build-seed:
	mkdir -p ${BUILD_DIR}
	cd ${BUILD_DIR}
	go build -mod vendor -o ${BUILD_DIR}/${SEED_BINARY_NAME} ./cmd/seed/main.go

.PHONY: run
run:
	cd ${BUILD_DIR} && ./${BINARY_NAME}

.PHONY: seed
seed:
	cd ${BUILD_DIR} && ./${SEED_BINARY_NAME}

.PHONY: copy-configs
copy-configs:
	mkdir -p ${BUILD_DIR}
	for path in `ls -d ./configs/* | sed 's:^\./::'` ; do \
  		mkdir -p ${BUILD_DIR}/$$path; \
		cp ./$$path/config.env ${BUILD_DIR}/$$path/config.env; \
	done

.PHONY: clean
clean:
	go clean
	rm -rf ./out

.PHONY: vendor
vendor:
	go mod vendor

.PHONY: deps
deps:
	go get ./...

.PHONY: vet
vet:
	go vet

.PHONY: lint
lint:
	golangci-lint run

.PHONY: lint-fix
lint-fix:
	golangci-lint run --fix

#	Docker section
.PHONY:	docker-compose-start-service
docker-compose-start-service:
	docker-compose -f ${DOCKER_COMPOSE_FILE} -p ${BINARY_NAME} up -d ${SERVICE_NAME}

.PHONY:docker-compose-stop-service
docker-compose-stop-service:
	docker-compose -f ${DOCKER_COMPOSE_FILE} stop ${SERVICE_NAME}
	docker-compose -f ${DOCKER_COMPOSE_FILE} rm --force ${SERVICE_NAME}

#	Database section
.PHONY:start-database
start-database:
	SERVICE_NAME=postgres $(MAKE) docker-compose-start-service

.PHONY:stop-database
stop-database:
	SERVICE_NAME=postgres $(MAKE) docker-compose-stop-service

.PHONY:init-database
init-database:
	./scripts/database/init.sh

.PHONY: migrate-up
migrate-up:
	./scripts/database/migrate-up.sh

.PHONY: migrate-down
migrate-down:
	./scripts/database/migrate-down.sh

#	Redis section
.PHONY:start-redis
start-redis:
	SERVICE_NAME=redis $(MAKE) docker-compose-start-service

.PHONY:stop-redis
stop-redis:
	SERVICE_NAME=redis $(MAKE) docker-compose-stop-service
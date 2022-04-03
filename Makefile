BUILD_DIR=./out/bin
BINARY_NAME=cryptomath-captcha
SEED_BINARY_NAME=${BINARY_NAME}-seed
MIGRATE_BINARY_NAME=${BINARY_NAME}-migrate

#	Go section
.PHONY: build-app
build-app:
	mkdir -p ${BUILD_DIR}
	cd ${BUILD_DIR}
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 GO111MODULE=on go build -v -mod vendor -o ${BUILD_DIR}/${OUT_NAME} ./cmd/${APP_NAME}/main.go

.PHONY: build-seed
build-seed:
	OUT_NAME=${BINARY_NAME} APP_NAME=seed $(MAKE) build-app

.PHONY: build-migrate
build-migrate:
	OUT_NAME=${MIGRATE_BINARY_NAME} APP_NAME=migrate $(MAKE) build-app

.PHONY: clean
clean:
	go clean
	rm -rf ${BUILD_DIR}

.PHONY: vendor
vendor:
	go mod vendor

.PHONY: install-deps
deps:
	go mod download

.PHONY: vet
vet:
	go vet

.PHONY: lint
lint:
	golangci-lint run

.PHONY: lint-fix
lint-fix:
	golangci-lint run --fix

#	Database section
.PHONY: migrate-up
migrate-up:
	migrate -path ./migrations -database "$(DATABASE_URL)" -verbose up

.PHONY: migrate-down
migrate-down:
	migrate -path ./migrations -database "$(DATABASE_URL)" -verbose down

.PHONE: boil-generate
boil-generate:
	sqlboiler psql
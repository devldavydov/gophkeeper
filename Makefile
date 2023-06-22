AGENT_VERSION := 1.0.0
SERVER_VERSION := 1.0.0
BUILD_DATE := $(shell date +'%d.%m.%Y %H:%M:%S')
BUILD_COMMIT := $(shell git rev-parse --short HEAD)

.PHONY: all
all: clean build lint test

.PHONY: build
build: build_client build_server

.PHONY: build_client
build_client:
	@echo "\n### $@"
	@mkdir -p ./bin
	@cd cmd/client && \
	 go build \
	 -ldflags "-X main.buildVersion=$(AGENT_VERSION) -X 'main.buildDate=$(BUILD_DATE)' -X main.buildCommit=$(BUILD_COMMIT)" \
	 -o ../../bin/client .

.PHONY: build_server
build_server:
	@echo "\n### $@"
	@mkdir -p ./bin
	@cd cmd/server && \
	 go build \
	 -ldflags "-X main.buildVersion=$(SERVER_VERSION) -X 'main.buildDate=$(BUILD_DATE)' -X main.buildCommit=$(BUILD_COMMIT)" \
	 -o ../../bin/server .

.PHONY: test
test:
	@echo "\n### $@"
	@echo "DON'T FORGET TO START postgres.sh\n"
	@export TEST_DATABASE_DSN=postgres://postgres:postgres@127.0.0.1:5432/praktikum?sslmode=disable && \
	 go test ./... -v --count 1

.PHONY: test_cover
test_cover:
	@echo "\n### $@"
	@echo "DON'T FORGET TO START postgres.sh\n"
	@export TEST_DATABASE_DSN=postgres://postgres:postgres@127.0.0.1:5432/praktikum?sslmode=disable && \
	 go test ./... -coverprofile cover.html -v --count 1
	@go tool cover -html=cover.html

.PHONY: run_docs
run_docs:
	@echo "See docs in http://localhost:8080/pkg/github.com/devldavydov/gophkeeper?m=all"
	@godoc -http=:8080

.PHONY: gen_tls
gen_tls:
	@echo "\n### $@"
	@mkdir -p tls
	@rm -rf tls/*
	@echo "subjectAltName=IP:127.0.0.1" > tls/server-ext.cnf
	@echo "> Generate CA's private key and self-signed certificate"
	@openssl req -x509 -newkey rsa:4096 -days 365 -nodes -keyout tls/ca-key.pem -out tls/ca-cert.pem -subj "/C=RU/ST=Moscow/L=Moscow/O=Yandex/OU=Praktikum/CN=gophkeeper"
	@echo "> CA's self-signed certificate"
	@openssl x509 -in tls/ca-cert.pem -noout -text
	@echo "> Generate web server's private key and certificate signing request (CSR)"
	@openssl req -newkey rsa:4096 -nodes -keyout tls/server-key.pem -out tls/server-req.pem -subj "/C=RU/ST=Moscow/L=Moscow/O=Yandex/OU=Praktikum/CN=gophkeeper"
	@echo "> Use CA's private key to sign web server's CSR and get back the signed certificate"
	@openssl x509 -req -in tls/server-req.pem -days 365 -CA tls/ca-cert.pem -CAkey tls/ca-key.pem -CAcreateserial -out tls/server-cert.pem -extfile tls/server-ext.cnf
	@echo "> Server's signed certificate"
	@openssl x509 -in tls/server-cert.pem -noout -text

.PHONY: gen_proto
gen_proto:
	@echo "\n### $@"
	@protoc --go_out=. --go_opt=paths=import --go-grpc_out=. --go-grpc_opt=paths=import internal/grpc/proto/user.proto
	@protoc --go_out=. --go_opt=paths=import --go-grpc_out=. --go-grpc_opt=paths=import internal/grpc/proto/secret.proto
	@protoc --go_out=. --go_opt=paths=import --go-grpc_out=. --go-grpc_opt=paths=import internal/grpc/proto/service.proto

.PHONY: gen_mock
gen_mock:
	@echo "\n### $@"
	@mockgen -destination=internal/grpc/mocks/mock_service_client.go -package=mocks github.com/devldavydov/gophkeeper/internal/grpc GophKeeperServiceClient

.PHONY: clean
clean:
	@echo "\n### $@"
	@rm -rf ./bin		

.PHONE: lint
lint:
	@echo "\n### $@"
	@golangci-lint --new-from-rev main run ./...
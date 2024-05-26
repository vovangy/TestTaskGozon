OS := $(shell uname -s)

ifneq ("$(wildcard .env)","")
include .env
endif

ifeq ($(OS), Linux)
	DOCKER_COMPOSE := docker compose
endif
ifeq ($(OS), Darwin)
	DOCKER_COMPOSE := docker-compose
endif

run:
	$(DOCKER_COMPOSE) up -d --build

stop:
	$(DOCKER_COMPOSE) down

test:
	go test -race ./...

coverage:
	go test -json ./... -coverprofile coverprofile_.tmp -coverpkg=./... ; \
	cat coverprofile_.tmp |grep -v auth.go| grep -v interfaces.go | grep -v docs.go| grep -v cors.go| grep -v transaction.go| grep -v main.go > coverprofile.tmp ; \
	rm coverprofile_.tmp ; \
	go tool cover -html coverprofile.tmp ; \
	go tool cover -func coverprofile.tmp

proto_gen:
	protoc -I proto proto/*.proto --go_out=./ --go-grpc_out=./


# If the first argument is "run"...
ifeq (migration,$(firstword $(MAKECMDGOALS)))
  # use the rest as arguments for "run"
  RUN_ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
  # ...and turn them into do-nothing targets
  $(eval $(RUN_ARGS):;@:)
endif

export GOBIN := $(PWD)/bin
export PATH := $(GOBIN):$(PATH)
export GO111MODULE := on
export PROTOC_VERSION := 3.9.1
export PROTOC_ZIP := protoc-$(PROTOC_VERSION)-linux-x86_64.zip
export PROTOC_VALIDATE_VERSION := 0.1.0
export PROTOC_VALIDATE := v$(PROTOC_VALIDATE_VERSION).zip

protoc:
	# protoc install
	curl -OL https://github.com/google/protobuf/releases/download/v$(PROTOC_VERSION)/$(PROTOC_ZIP)
	unzip -o $(PROTOC_ZIP) -d ./ bin/protoc
	unzip -o $(PROTOC_ZIP) -d ./proto include/*
	rm -f $(PROTOC_ZIP)

# protoc-gen-validate installation
protoc_gen_validate:
	curl -OL https://github.com/envoyproxy/protoc-gen-validate/archive/v$(PROTOC_VALIDATE_VERSION).zip
	unzip -j $(PROTOC_VALIDATE) protoc-gen-validate-$(PROTOC_VALIDATE_VERSION)/validate/* -d ./proto/include/validate
	rm -f $(PROTOC_VALIDATE)

tools: protoc
	# required go tools installation
	go install github.com/gojuno/minimock/v3/cmd/minimock
	go install github.com/hexdigest/gowrap/cmd/gowrap
	go install github.com/envoyproxy/protoc-gen-validate
	go install github.com/twitchtv/twirp/protoc-gen-twirp
	go install github.com/golangci/golangci-lint/cmd/golangci-lint
	go install github.com/golang/protobuf/protoc-gen-go
	go install github.com/gojuno/goose/cmd/goose

generate:
	go generate
	bash .werffiles/create_config.sh | tee config.local.yml

lint:
	golangci-lint run ./... --skip-files ".*(_mock_|_with_fallback|_with_prometheus|_with_tracing|_with_error_logging|_with_validation).*.go$\"

test: lint
	GOGC=off go test -race ./...

build: migration_tools
	go build -o ./bin/app ./cmd/app

run: generate build
	./bin/app -f config.local.yml

# use https://hub.docker.com/r/erichsu/redis-sentinel/ to simulate redis cluster
infrastructure:
	docker-compose -f docker-compose.yml -p sirocco start

migration_tools:
	go install github.com/gojuno/goose/cmd/goose

migration:
	goose -dir ./storage/migrations create $(RUN_ARGS) sql

migrate_up:
	goose -dir ./storage/migrations -conf config.local.yml up

migrate_down:
	goose -dir ./storage/migrations -conf config.local.yml down

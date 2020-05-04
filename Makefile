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

lint:
	golangci-lint run ./... --skip-files ".*(_mock_|_with_fallback|_with_prometheus|_with_tracing|_with_error_logging|_with_validation).*.go$\"

test: lint
	GOGC=off go test -race ./...

build: migration_tools
	go build -o ./bin/app ./cmd/app


migration_tools:
	go install github.com/gojuno/goose/cmd/goose

migration:
	goose -dir ./storage/migrations create $(RUN_ARGS) sql

#migrate_up:
#	goose -dir ./storage/migrations -conf config.local.yml up
#
#migrate_down:
#	goose -dir ./storage/migrations -conf config.local.yml down

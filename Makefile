
.DEFAULT_GOAL := help

GOBINDATA=$(shell type -p go-bindata 2>/dev/null)
VERSION:=$(shell git describe --always --long)
PROJECT_NAME:= vault_cassandra
CLONE_URL:=github.com/daiemna/$(PROJECT_NAME)
IDENTIFIER= $(VERSION)-$(GOOS)-$(GOARCH)
BUILD_TIME=$(shell date -u +%FT%T%z)
LDFLAGS='-extldflags "-static" -s -w -X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME)'

.PHONY : test environment
help:          ## Show available options with this Makefile
	@grep -F -h "##" $(MAKEFILE_LIST) | grep -v grep | awk 'BEGIN { FS = ":.*?##" }; { printf "%-18s  %s\n", $$1,$$2 }'

.PHONY: test
test:          ## Run all the tests
	chmod +x ./scripts/test.sh && ./scripts/test.sh


.PHONY: environment
environment:  ## Recreate the test environment
	docker-compose -f environment/docker-compose.yml down --remove-orphans && \
	docker-compose -f environment/docker-compose.yml up -d 
	environment/scripts/wait-for-it.sh -t 40 localhost:9042 -- echo "cassandra active, will wait for 60s to bootstrap"
	sleep 60
	docker exec -it int-test-scylla-c cqlsh -u cassandra -p cassandra -f /etc/init.cql 
	environment/scripts/vault_init.sh

.PHONY: run
run:         ## Clean the application
	@go clean -i ./...
	@rm -rf ./$(PROJECT_NAME)
	@rm -rf build/*

.PHONY: vendor
vendor:           ## Go get vendor
	go mod vendor
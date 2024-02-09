BASH_PATH:=$(shell which bash)
SHELL=$(BASH_PATH)
ROOT := $(shell realpath $(dir $(lastword $(MAKEFILE_LIST))))
CURRENT_TIMESTAMP := $(shell date +%s)

all: format lint-ci test

prepare-config:
	cp config.example.yaml config.yaml

build-docker:
	sudo docker build . -f ./Dockerfile -t chuki:latest

prepare-compose:
	mkdir .compose
	cp ./config.example.yaml ./.compose/config.yaml

check-compose:
	@if [[ ! -f .compose/config.yaml ]]; then \
		echo ".compose/config.yaml doesnt exists"; \
		echo "run `make prepare-compose`"; \
		exit 1; \
	fi;

docker:
	sudo docker build -t chuki .

up: check-compose
	sudo docker-compose up -d

down:
	sudo docker-compose down

full-down:
	sudo docker-compose -f ./deployments/docker-compose.yaml down --remove-orphans --volumes

build-me:
	go build -v -o chuki -race .

run: build-me prepare-config
	./chuki serve --config config.yaml

db-migrate: build-me
	./chuki migrate --config config.yaml

############################################################
# Format & Lint
############################################################
check-goimport:
	which goimports || GO111MODULE=off go get -u golang.org/x/tools/cmd/goimports

format: check-goimport
	find $(ROOT) -type f -name "*.go" -not -path "$(ROOT)/vendor/*" | xargs -I R -n 1 goimports -w R
	find $(ROOT) -type f -name "*.go" -not -path "$(ROOT)/vendor/*" | xargs -I R -n 1 gofmt -s -w R

lint:
	find $(ROOT) -type f -name "*.go" -not -path "$(ROOT)/vendor/*" | xargs -I R -n 1 echo R | xargs -I R -n 1 go vet R

check-golangci-lint:
	which golangci-lint || (curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sudo sh -s -- -b $(go env GOPATH)/bin v1.55.0)

lint-ci: check-golangci-lint vendor
	golangci-lint run -c .golangci.yml ./...

vendor:
	go mod vendor -v

############################################################
# Test & Coverage
############################################################
check-gotestsum:
	which gotestsum || (go install gotest.tools/gotestsum@latest)

check-gocover-cobertura:
	which gocover-cobertura || (go install github.com/boumenot/gocover-cobertura@latest)

test: check-gotestsum vendor
	gotestsum --junitfile-testcase-classname short --junitfile .report.xml -- -gcflags 'all=-N -l' -mod vendor ./...

coverage: check-gotestsum vendor
	gotestsum -- -gcflags 'all=-N -l' -mod vendor -v -coverprofile=.testCoverage.txt ./...
	GOFLAGS=-mod=vendor go tool cover -func=.testCoverage.txt

coverage-report: check-gocover-cobertura coverage
	GOFLAGS=-mod=vendor go tool cover -html=.testCoverage.txt -o testCoverageReport.html
	gocover-cobertura < .testCoverage.txt >  .cobertura.raw.xml
	sed 's;filename=\"github.com/Shopitant/backend/;filename=\";g' .cobertura.raw.xml > .cobertura.xml

PROJECT_DIR=$(dir $(realpath $(firstword $(MAKEFILE_LIST))))

BINARY_NAME = aryzona
DOCKER_IMAGE_NAME = aryzonabot

DIST_LDFLAGS = $(LDFLAGS) -w -s
TEST_COMMAND=go test

# kinda shitty and hacky what of doing that... =(
COMMIT_MESSAGE = $(shell git log -1 --pretty=%s | sed "s/'//g; s/\"//g")
COMMIT_HASH = $(shell git rev-list -1 HEAD)

LDFLAGS = -X 'main.commitMessage=$(COMMIT_MESSAGE)' -X 'main.commitHash=$(COMMIT_HASH)'

.PHONY: build
build:
	CGO_ENABLED=0 go build -v -ldflags="$(LDFLAGS)" -o $(PROJECT_DIR)$(BINARY_NAME) $(PROJECT_DIR)cmd/aryzona

.PHONY: run
run: build
	./$(BINARY_NAME) 

.PHONY: docker
docker:
	docker build -t $(DOCKER_IMAGE_NAME) .

.PHONY: install
install: build
	sudo cp ./$(BINARY_NAME) /usr/bin/

.PHONY: test
test: 
	$(TEST_COMMAND) -cover -parallel 5 -failfast  ./... 

.PHONY: tidy
tidy:
	go mod tidy

# (build but with a smaller binary)
.PHONY: dist
dist:
	CGO_ENABLED=0 go build -gcflags=all=-l -v -ldflags="$(DIST_LDFLAGS)" -o $(PROJECT_DIR)$(BINARY_NAME) $(PROJECT_DIR)cmd/aryzona

# (even smaller binary)
.PHONY: pack
pack: dist
	upx ./$(BINARY_NAME)

.PHONY: lint
lint:
	revive -formatter friendly -config revive.toml ./...

.PHONY: spell
spell:
	misspell -error ./**

.PHONY: staticcheck
staticcheck:
	staticcheck ./...

.PHONY: gosec
gosec:
	gosec -tests ./... 

.PHONY: inspect
inspect: lint spell gosec staticcheck

# auto restart bot (using fiber CLI)
.PHONY: dev
dev:
	fiber dev

BINARY_NAME = aryzona
DOCKER_IMAGE_NAME = aryzonabot

DIST_LDFLAGS = $(LDFLAGS) -w -s
TEST_COMMAND=go test

# FIXME: kinda shitty and hacky way of doing that... =(
COMMIT_MESSAGE = $(shell git log -1 --pretty=%s | sed "s/'//g; s/\"//g")
COMMIT_HASH = $(shell git rev-list -1 HEAD)

LDFLAGS = -X 'main.commitMessage=$(COMMIT_MESSAGE)' -X 'main.commitHash=$(COMMIT_HASH)'

.PHONY: build
build:
	CGO_ENABLED=0 go build -v -ldflags="$(LDFLAGS)" -o $(BINARY_NAME) ./cmd/$(BINARY_NAME)

.PHONY: run
run: build
	./$(BINARY_NAME) 

.PHONY: docker
docker:
	docker build -t $(DOCKER_IMAGE_NAME) .

.PHONY: install
install: build
	sudo cp ./$(BINARY_NAME) /usr/local/bin/

.PHONY: test
test: 
	$(TEST_COMMAND) -cover -parallel 5 -failfast -count=1 ./... 

# human readable test output
.PHONY: love
love:
ifeq ($(filter watch,$(MAKECMDGOALS)),watch)
	gotestsum --format-hide-empty-pkg --watch ./...
else
	gotestsum --format-hide-empty-pkg ./...
endif

.PHONY: tidy
tidy:
	go mod tidy

# (build but with a smaller binary)
.PHONY: dist
dist:
	CGO_ENABLED=0 go build -gcflags=all=-l -v -ldflags="$(DIST_LDFLAGS)" -o $(BINARY_NAME) ./cmd/$(BINARY_NAME)

# (even smaller binary)
.PHONY: pack
pack: dist
	upx ./$(BINARY_NAME)

.PHONY: dev
dev:
	air

.PHONY: lint
lint:
	revive -formatter friendly -config revive.toml ./...

.PHONY: staticcheck
staticcheck:
	staticcheck ./...

.PHONY: gosec
gosec:
	gosec -tests ./... 

.PHONY: inspect
inspect: lint gosec staticcheck

.PHONY: install-inspect-tools
install-inspect-tools:
	go install github.com/mgechev/revive@latest
	go install honnef.co/go/tools/cmd/staticcheck@latest
	go install github.com/securego/gosec/v2/cmd/gosec@latest

.PHONY: install-dev-tools
install-dev-tools: install-inspect-tools
	go install github.com/air-verse/air@latest
	go install gotest.tools/gotestsum@latest
	go install github.com/rubenv/sql-migrate/...@latest

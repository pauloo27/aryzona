BINARY_NAME = aryzona
COMMIT_MESSAGE = $(shell git log -1 --pretty=%s | sed "s/'//g; s/\"//g")
COMMIT_HASH = $(shell git rev-list -1 HEAD)
LDFLAGS = -X 'main.commitMessage=$(COMMIT_MESSAGE)' -X 'main.commitHash=$(COMMIT_HASH)'
DIST_LDFLAGS = $(LDFLAGS) -w -s
TEST_COMMAND=go test

.PHONY: build
build:
	go build -v -ldflags="$(LDFLAGS)" -o $(BINARY_NAME) ./cmd/aryzona

.PHONY: run
run: build
	./$(BINARY_NAME) 

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
	go build -gcflags=all=-l -v -ldflags="$(DIST_LDFLAGS)"

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

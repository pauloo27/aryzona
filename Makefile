BINARY_NAME = aryzona
COMMIT_MESSAGE = $(shell git log -1 --pretty=%s | sed "s/'//g; s/\"//g")
COMMIT_HASH = $(shell git rev-list -1 HEAD)
LDFLAGS = -X 'main.commitMessage=$(COMMIT_MESSAGE)' -X 'main.commitHash=$(COMMIT_HASH)'
DIST_LDFLAGS = $(LDFLAGS) -w -s
TEST_COMMAND=go test

build:
	go build -v -ldflags="$(LDFLAGS)"

run: build
	./$(BINARY_NAME) 

install: build
	sudo cp ./$(BINARY_NAME) /usr/bin/

test: 
	$(TEST_COMMAND) -cover -parallel 5 -failfast  ./... 

colorful_test: 
	gotest -cover -parallel 5 -failfast  ./... 

nocolor_test: 
	go test -cover -parallel 5 -failfast  ./... 

tidy:
	go mod tidy

# (build but with a smaller binary)
dist:
	go build -gcflags=all=-l -v -ldflags="$(DIST_LDFLAGS)"

# (even smaller binary)
pack: dist
	upx ./$(BINARY_NAME)

# kill previous version and start a new one 
restart_bot: build
	- killall $(BINARY_NAME) -w
	./$(BINARY_NAME) 

lint:
	revive -formatter friendly -config revive.toml ./...

spell:
	misspell -error ./**

staticcheck:
	staticcheck ./...

gosec:
	gosec -tests ./... 

inspect: lint spell gosec staticcheck

# auto restart bot (using fiber CLI)
dev:
	fiber dev

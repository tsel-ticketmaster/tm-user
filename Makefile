.PHONY: install test-dev test cover run.dev build clean

install:
	go mod download

test:
	@echo "Run unit testing ..."
		mkdir -p ./coverage && \
			go test -v -coverprofile=./coverage/coverage.out -covermode=atomic ./...

cover: test
	@echo "Generating coverprofile ..."
		go tool cover -func=./coverage/coverage.out &&\
			go tool cover -html=./coverage/coverage.out -o ./coverage/coverage.html

run.web:
	@echo "Run in development mode ..."
		go run cmd/web/main.go

build.web:
	@echo "Building the executable file ..."
		CGO_ENABLED=0 GOOS=linux go build -a -o bin/webapp cmd/web/main.go &&\
			cp bin/webapp /tmp/webapp

clean:
	@echo "Cleansing the last built ..."
		rm -rf bin
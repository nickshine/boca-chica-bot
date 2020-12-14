BINARY=boca-chica-bot
BUILD_FLAGS=-ldflags="-s -w"

lint:
	golangci-lint run -v

build:
	GOOS=linux GOARCH=amd64 go build $(BUILD_FLAGS) -o bin/$(BINARY) ./cmd

run: build
	docker run --env-file .env --rm -v ${PWD}/bin:/var/task:ro,delegated lambci/lambda:go1.x $(BINARY)

zip:
	@echo "Zipping for release"
	zip -j deploy/lambda/lambda.zip bin/$(BINARY)

package: clean build zip

clean:
	@go clean
	@rm -rf bin
	@rm -f deploy/lambda/*.zip
	@rm -f coverage.out coverage.html

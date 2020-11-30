BINARY=boca-chica-bot
BUILD_FLAGS=-ldflags="-s -w"

lint:
	golangci-lint run -v

build:
	GOOS=linux GOARCH=amd64 go build -o bin/$(BINARY) ./cmd

run: build
	docker run --env-file .env --rm -v ${PWD}/bin:/var/task:ro,delegated lambci/lambda:go1.x $(BINARY)

zip:
	@echo "Zipping for release"
	zip -jm bin/lambda.zip bin/$(BINARY)

package: clean build zip

clean:
	@go clean
	@rm -rf bin/
	@rm -f coverage.out coverage.html

BINARY=boca-chica-bot
BUILD_FLAGS=-ldflags="-s -w"

lint:
	golangci-lint run -v --timeout 5m

build:
	GOOS=linux GOARCH=amd64 go build -o bin/$(BINARY) ./lambda

zip:
	@echo "Zipping for release"
	@cd bin && zip lambda.zip $(BINARY); cd -

package: clean build zip

clean:
	@go clean
	@rm -rf bin/
	@rm -f coverage.out coverage.html

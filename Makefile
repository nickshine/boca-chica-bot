BINARY=boca-chica-bot
SCRAPER_BINARY=$(BINARY)-scraper
PUBLISHER_BINARY=$(BINARY)-publisher
BUILD_FLAGS=-ldflags="-s -w"

lint:
	golangci-lint run -v

test:
	go test -v -covermode=count -coverprofile=coverage.out ./...

cover:
	go tool cover -html=coverage.out -o coverage.html

build:
	GOOS=linux GOARCH=amd64 go build $(BUILD_FLAGS) -o bin/$(SCRAPER_BINARY) ./cmd/scraper
	GOOS=linux GOARCH=amd64 go build $(BUILD_FLAGS) -o bin/$(PUBLISHER_BINARY) ./cmd/publisher

payload:
	sam local generate-event dynamodb update > payload.json

run-scraper:
	docker run --env-file .env --rm -v ${PWD}/bin:/var/task:ro,delegated lambci/lambda:go1.x $(SCRAPER_BINARY)

run-publisher:
	docker run --env-file .env --rm -v ${PWD}/bin:/var/task:ro,delegated lambci/lambda:go1.x $(PUBLISHER_BINARY) '$(shell cat payload.json)'

zip:
	@echo "Zipping for release"
	zip -j deploy/lambda/scraper-lambda.zip bin/$(SCRAPER_BINARY)
	zip -j deploy/lambda/publisher-lambda.zip bin/$(PUBLISHER_BINARY)

package: clean build zip

clean:
	@go clean
	@rm -rf bin dist payload.json
	@rm -f deploy/lambda/*.zip
	@rm -f coverage.out coverage.html

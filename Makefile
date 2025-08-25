CURRENT_UID := $(shell id -u):$(shell id -g)

setup:
	go mod vendor

lint:
	docker run --rm -v $(CURDIR):/app -v ~/.cache/golangci-lint/v1.63.4:/root/.cache -w /app golangci/golangci-lint:v1.63.4 golangci-lint run ./...

test:
	go test -race ./...

fmt:
	gofmt -s -w -e .

generate-graphql:
	GO111MODULE=on go run -mod=mod github.com/Khan/genqlient

validate: fmt lint

build:
	go build -o nexaa .
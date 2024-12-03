NOW := $(shell date +%Y%m%d-%H%M%S)
test:
	https_proxy=127.0.0.1:8082 go run -tags dev . container create \
		--namespace 3093 \
		--name "test-$(NOW)" \
		--image nginx:latest \
		--port "80:80" \
		--https "hello.tilaa.dev"

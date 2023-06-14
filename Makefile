.PHONY: all
all:
	cd cmd/yagisan && go build

.PHONY: install-tools
install-tools:
	go install github.com/sqs/goreturns

.PHONY: format
format:
	goreturns -w .
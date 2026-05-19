export GO ?= go
export CGO_ENABLED = 0

TAG := $(shell git describe --always --tags $(git rev-list --tags --max-count=1) --match v*)

.PHONY: all
all: awgproxy

.PHONY: awgproxy
awgproxy:
	${GO} build -trimpath -ldflags "-s -w -X 'main.version=${TAG}'" ./cmd/awgproxy

.PHONY: clean
clean:
	${RM} awgproxy

.PHONY: codegen fix fmt vet lint test build tidy run

GOBIN := $(shell go env GOPATH)/bin

build:
	go build -o example-apiserver .

all: clean codegen fmt test build tidy

docker:
	GOOS=linux GOARCH=amd64 go build -o install/bin/apiserver
	docker build --tag IMAGE:VERSION install

install: docker
	kustomize build install | kubectl apply -f -

apiserver-logs:
	kubectl logs -l apiserver=true --container apiserver -n example-system -f

clean:
	find . -name "zz_generated.*" | xargs rm

codegen:
	./hack/update-codegen.sh

fix:
	go fix ./...

fmt:
	go fmt ./...

tidy:
	go mod tidy

lint:
	(which golangci-lint || go get github.com/golangci/golangci-lint/cmd/golangci-lint)
	$(GOBIN)/golangci-lint run ./...

test:
	go test -cover ./...

vet:
	go vet ./...

run:
	./hack/run.sh

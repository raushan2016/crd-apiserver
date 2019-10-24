
# Image URL to use all building/pushing image targets
IMG ?= raushan2016/crd-apiserver1:latest
OUTPUT ?= ./bin/final.yaml
# Produce CRDs that work back to Kubernetes 1.11 (no version conversion)
CRD_OPTIONS ?= "crd:trivialVersions=true"
DEV_CONTROLLER_ID ?= $(shell whoami)
ENV ?= dataplane-dev
export GO111MODULE=on

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

mod:
	go mod tidy

# Build apiserver binary
manager: mod fmt vet
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o bin/apiserver main.go

generate:
	kustomize build config -o ${OUTPUT}

# Run against the configured Kubernetes cluster in ~/.kube/config
run: fmt vet manifests
	go run ./main.go

# Deploy controller in the configured Kubernetes cluster in ~/.kube/config
deploy:
	kubectl apply -f ${OUTPUT}

dev-deploy: manifests docker-build docker-push deploy

# Deletes controller in the configured Kubernetes cluster in ~/.kube/config
undeploy:
	kubectl delete -f ${OUTPUT}

# Run go fmt against code
fmt:
	go fmt ./...

# Run go vet against code
vet:
	go vet ./...

docker-all: docker-build docker-push

# Build the docker image
docker-build: manager
	docker build . -t ${IMG}

# Push the docker image
docker-push:
	docker push ${IMG}

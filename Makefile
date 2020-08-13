SHELL := /bin/bash
BASEDIR = $(shell pwd)
export GO111MODULE=on
export GOPROXY=https://goproxy.cn,direct
export GOPRIVATE=*.gitlab.com
export GOSUMDB=off

APP_NAME=`cat package.json | grep name | head -1 | awk -F: '{ print $$2 }' | sed 's/[\",]//g' | tr -d '[[:space:]]'`
APP_VERSION=`cat package.json | grep version | head -1 | awk -F: '{ print $$2 }' | sed 's/[\",]//g' | tr -d '[[:space:]]'`
COMMIT_ID=`git rev-parse HEAD`
IMAGE_PREFIX="registry.cn-hangzhou.aliyuncs.com/makeblock/${APP_NAME}:v${APP_VERSION}"

fmt:
	gofmt -w .
utest:
	go mod tidy; \
	go test -coverpkg=./... -coverprofile=coverage.data ./...;
lint:
	golangci-lint run -c .golangci.yml;
build:
	go mod tidy; \
	cd deploy/docker;  \
	rm -rf build | echo "no build dir"; \
	mkdir build; \
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags '-w -s' -o build/main ../../main.go;
build-master: build
	IMAGE_NAME="${IMAGE_PREFIX}-master"; \
	echo $$IMAGE_NAME; \
	cd deploy/docker;  \
	docker build --build-arg tmp_api_version=${COMMIT_ID} -t $$IMAGE_NAME -f Dockerfile .; \
	rm -rf build | echo "no build dir"; \
	docker push $$IMAGE_NAME;
build-release: build
	IMAGE_NAME="${IMAGE_PREFIX}-release"; \
	echo $$IMAGE_NAME; \
	cd deploy/docker;  \
	docker build --build-arg tmp_api_version=${COMMIT_ID} -t $$IMAGE_NAME -f Dockerfile .; \
	rm -rf build | echo "no build dir"; \
	docker push $$IMAGE_NAME;
deploy-dev:
	NEW_IMAGE="${APP_NAME}=${IMAGE_PREFIX}-master"; \
	cd deploy/kubernetes/base; \
	kustomize edit set image $$NEW_IMAGE; \
	kustomize edit add annotation git.commit.id:${COMMIT_ID}
	kubectl config use aliyun-test
	kustomize build deploy/kubernetes/overlays/dev
	kustomize build deploy/kubernetes/overlays/dev | kubectl apply -f -
deploy-test:
	NEW_IMAGE="${APP_NAME}=${IMAGE_PREFIX}-release"; \
	cd deploy/kubernetes/base; \
	kustomize edit set image $$NEW_IMAGE; \
	kustomize edit add annotation git.commit.id:${COMMIT_ID}
	kubectl config use aliyun-test
	kustomize build deploy/kubernetes/overlays/test
	kustomize build deploy/kubernetes/overlays/test | kubectl apply -f -
deploy-we-test:
	NEW_IMAGE="${APP_NAME}=${IMAGE_PREFIX}-release"; \
	cd deploy/kubernetes/base; \
	kustomize edit set image $$NEW_IMAGE; \
	kustomize edit add annotation git.commit.id:${COMMIT_ID}
	kubectl config use azure
	kustomize build deploy/kubernetes/overlays/we-test
	kustomize build deploy/kubernetes/overlays/we-test | kubectl apply -f -
deploy-pre:
	NEW_IMAGE="${APP_NAME}=${IMAGE_PREFIX}-release"; \
	cd deploy/kubernetes/base; \
	kustomize edit set image $$NEW_IMAGE; \
	kustomize edit add annotation git.commit.id:${COMMIT_ID}
	kubectl config use aliyun
	kustomize build deploy/kubernetes/overlays/pre
	kustomize build deploy/kubernetes/overlays/pre | kubectl apply -f -
deploy-prod:
	kubectl config use aliyun
	kustomize build deploy/kubernetes/overlays/prod | kubectl apply -f -
deploy-prod-preview:
	NEW_IMAGE="${APP_NAME}=${IMAGE_PREFIX}-release"; \
	cd deploy/kubernetes/base; \
	kustomize edit set image $$NEW_IMAGE; \
	kustomize edit add annotation git.commit.id:${COMMIT_ID};
	kustomize build deploy/kubernetes/overlays/prod
deploy-we-prod:
	kubectl config use azure
	kustomize build deploy/kubernetes/overlays/we-prod | kubectl apply -f -
deploy-we-prod-preview:
	NEW_IMAGE="${APP_NAME}=${IMAGE_PREFIX}-release"; \
	cd deploy/kubernetes/base; \
	kustomize edit set image $$NEW_IMAGE; \
	kustomize edit add annotation git.commit.id:${COMMIT_ID};
	kustomize build deploy/kubernetes/overlays/we-prod	
help:
	@echo "fmt - format the source code"
	@echo "utest - unit test"
	@echo "build - build amd64 binary"
	@echo "build-master - build for master branch"
	@echo "build-release - build for release branch"
	@echo "deploy-dev - deploy to shenzhen develop environment"
	@echo "deploy-test - deploy to shenzhen test environment"
	@echo "deploy-we-test - deploy to west europe test environment"
	@echo "deploy-pre - deploy to shenzhen preview environment"
	@echo "deploy-prod - deploy to shenzhen production environment"
	@echo "deploy-prod-preview - preview the yaml of shenzhen production environment"
	@echo "deploy-we-prod - deploy to west europe production environment"
	@echo "deploy-we-prod-preview - preview the yaml of west europe production environment"
export GO111MODULE=on
PROJECT=perindo
NAME=be-service-auth
TAG := $(shell git describe --candidates=0 2>/dev/null)
BRANCH := $(shell git rev-parse --abbrev-ref HEAD 2>/dev/null)
VERSION := $(TAG:v%=%)
ifeq ($(TAG),)
	VERSION := $(BRANCH)
endif

init:
	git config core.hooksPath .githooks
	go mod tidy

build:
	#cp api-specification/openapi.yaml saksi/delivery/http/openapi
	#cp api-specification/grpc.proto saksi/delivery/grpc/proto
	# rm saksi/delivery/grpc/auth/*
	# mkdir -p saksi/delivery/grpc/auth
	# protoc --go_out=saksi/delivery/grpc/auth --go_opt=paths=source_relative --go-grpc_out=saksi/delivery/grpc/auth --go-grpc_opt=paths=source_relative --proto_path=saksi/delivery/grpc/proto grpc.proto
	# protoc --plugin=protoc-gen-doc=${HOME}/go/bin/protoc-gen-doc --doc_out=saksi/delivery/http/grpcdoc --doc_opt=saksi/delivery/http/grpcdoc/html.tmpl,index.html --proto_path=saksi/delivery/http/grpcdoc grpc.proto
	go mod tidy
	go build -o ${NAME} app/*.go

clean:
	if [ -f saksi/delivery/http/openapi/openapi.yaml ] ; then rm saksi/delivery/http/openapi/openapi.yaml ; fi
	if [ -f saksi/delivery/http/grpcdoc/grpc.proto ] ; then rm saksi/delivery/http/grpcdoc/grpc.proto ; fi
	if [ -f ${NAME} ] ; then rm ${NAME} ; fi

docker:
	docker build -t ${REGISTRY}/${PROJECT}/${NAME}:$(VERSION) .

run:
	cp api-specification/openapi.yaml saksi/delivery/http/openapi
	cp api-specification/grpc.proto saksi/delivery/http/grpcdoc
	# rm saksi/delivery/grpc/auth/*
	# mkdir -p saksi/delivery/grpc/auth
	# protoc --go_out=saksi/delivery/grpc/auth --go_opt=paths=source_relative --go-grpc_out=saksi/delivery/grpc/auth --go-grpc_opt=paths=source_relative --proto_path=saksi/delivery/grpc/proto grpc.proto
	# protoc --plugin=protoc-gen-doc=${HOME}/go/bin/protoc-gen-doc --doc_out=saksi/delivery/http/grpcdoc --doc_opt=saksi/delivery/http/grpcdoc/html.tmpl,index.html --proto_path=saksi/delivery/http/grpcdoc grpc.proto
	go mod tidy
	go run app/*.go -c config.yaml

push:
	docker push ${REGISTRY}/${PROJECT}/${NAME}:$(VERSION)

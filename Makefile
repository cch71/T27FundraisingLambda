ifneq (,$(wildcard ./.env))
    include .env
    export
endif

MK_DIR := $(realpath $(dir $(lastword $(MAKEFILE_LIST))))
DIST_DIR := ${MK_DIR}/dist

clean:
	@rm -rf dist
	@mkdir -p dist

lambda: clean
	cd ${MK_DIR}/cmd/lambda && GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -o ${DIST_DIR}/gqlhandler

cli: clean
	cd ${MK_DIR}/cmd/t27frcli && go build -o ${DIST_DIR}/t27frcli

dist: lambda
	cp $(DB_CA_ROOT_PATH) dist
	cd dist && zip function.zip gqlhandler root.crt

run:
	aws-sam-local local start-api

install:
	go get github.com/aws/aws-lambda-go/events
	go get github.com/aws/aws-lambda-go/lambda
	go get github.com/stretchr/testify/assert

install-dev:
	go get github.com/awslabs/aws-sam-local

update-mods:
	cd frgql && go get -u ./... && go mod tidy
	cd cmd/lambda && go get -u ./... && go mod tidy
	cd cmd/t27frcli && go get -u ./... && go mod tidy


test:
	go test ./... --cover

deploy: dist
	op plugin run -- aws lambda update-function-code --function-name ${GQL_LAMBDA_FUNCTION_NAME} --zip-file fileb://${PWD}/dist/function.zip

syncusers:
	cd cmd/t27frcli && op run --env-file="../../.env" -- go run . syncusers


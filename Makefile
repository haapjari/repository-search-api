MAIN_MODULE             := cmd/main.go
IMAGE_VERSION           := latest
SERVICE_NAME 			?= repository-search-api
OUTPUT_PATH 			?= bin/repository-search-api

##  
## Commands
##

.PHONY: run
run: compile
	./${OUTPUT_PATH}

.PHONY: dev
dev:
	go run -race cmd/main.go

.PHONY: test
test:
	go test -v -count=1 ./...


.PHONY: docker-build
docker-build:
	docker build -t ${SERVICE_NAME}:${IMAGE_VERSION} .


.PHONY: compile
compile:
	go build -o ${OUTPUT_PATH} ${MAIN_MODULE}

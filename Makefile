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

docker-debug:
	docker build -t ${SERVICE_NAME}:debug -f Dockerfile.Debug .
.PHONY: docker-debug

docker-debug-run:
	docker run -itd -e PORT=8000 -p 8000:8000 --name ${SERVICE_NAME} ${SERVICE_NAME}:debug


.PHONY: compile
compile:
	go build -o ${OUTPUT_PATH} ${MAIN_MODULE}

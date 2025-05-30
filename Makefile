MAIN_MODULE             := cmd/main.go
IMAGE_VERSION           := latest
SERVICE_NAME 			?= repository-search-api
OUTPUT_PATH 			?= bin/repository-search-api

#  
## Commands
#


.PHONY: run
run:
	go run ${MAIN_MODULE}


.PHONY: run/dev
run/dev:
	go run -race ${MAIN_MODULE}


.PHONY: test
test:
	go test -v -count=1 ./...


.PHONY: docker/build
docker/build:
	docker build -t ${SERVICE_NAME}:${IMAGE_VERSION} .


.PHONY: docker/run
docker/run:
	docker run -itd -e PORT=8000 -p 8000:8000 --name ${SERVICE_NAME} ${SERVICE_NAME}



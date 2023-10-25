include make.properties

run: compile
	${OUTPUT_PATH}

workspace:
	go work use .

dev:
	air

test:
	go clean
	go test ./...

compile:
	go build -o ${OUTPUT_PATH} ${MAIN_MODULE}

docker:
	docker build --tag ${DOCKER_IMAGE}:latest .

docker-run:
	docker run -idt -p 8080:8080 --name ${DOCKER_IMAGE} ${DOCKER_IMAGE}:latest

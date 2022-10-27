include make.properties

run:
	go run ${MAIN_MODULE}

test:
	go clean
	go test ./...

compile:
	go build -o ${OUTPUT_PATH} ${MAIN_MODULE}

run-bin:
	${OUTPUT_PATH}

docker:
	docker build --tag ${DOCKER_IMAGE} .

docker-run:
	docker run -id -p 8080:8080 ${DOCKER_IMAGE}:${IMAGE_VERSION}

database-start:
	sudo service postgresql start

database-status:
	sudo service postgresql status

database-stop:
	sudo service postgresql stop

database-exec:
	sudo -u postgres psql

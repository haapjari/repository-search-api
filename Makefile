include make.properties

run:
	go run ${MAIN_MODULE}

test:
	go mod tidy
	mkdir out && go test -v -coverprofile out/cover.out ./...

build:
	go build -o ${OUTPUT_PATH} ${MAIN_MODULE}

run-bin:
	${OUTPUT_PATH}

build-docker:
	docker build --tag ${DOCKER_IMAGE} .

run-docker:
	docker run -id -p 8080:8080 ${DOCKER_IMAGE}:${IMAGE_VERSION}

database-start:
	sudo service postgresql start

database-status:
	sudo service postgresql status

database-stop:
	sudo service postgresql stop

database-exec:
	sudo -u postgres psql

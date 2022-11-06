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

docker-build:
	docker build --tag ${DOCKER_IMAGE}:latest .

docker-run:
	docker run -idt -p 8080:8080 --name ${DOCKER_IMAGE} --net glass_glass --ip ${DOCKER_STATIC_IP} ${DOCKER_IMAGE}:latest

docker-compose:
	docker-compose up -d

database-start:
	sudo service postgresql start

database-status:
	sudo service postgresql status

database-stop:
	sudo service postgresql stop

database-exec:
	sudo -u postgres psql

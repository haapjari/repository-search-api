include make.properties

run:
	make compile
	make run-bin

get-all:
	./requests/repository-get-all.sh

fetch-metadata:
	./requests/repository-fetch-metadata.sh

workspace:
	go work use .

air:
	air

test:
	go clean
	go test ./...

compile:
	go build -o ${OUTPUT_PATH} ${MAIN_MODULE}

profile-cpu:
	go run -cpuprofile cpu.prof ${MAIN_MODULE}

profile-memory:
	go run -memprofile mem.prof ${MAIN_MODULE}

profile-heap:
	go run -memprofilerate heap.prof ${MAIN_MODULE}

# profile-analyze:
# go tool pprof myprof cpu.prof

run-bin:
	${OUTPUT_PATH}

docker-build:
	docker build --tag ${DOCKER_IMAGE}:latest .

docker-run:
	docker run -idt -p 8080:8080 --name ${DOCKER_IMAGE} --net ${DOCKER_NETWORK} --ip ${DOCKER_STATIC_IP} ${DOCKER_IMAGE}:latest

docker-compose:
	docker-compose up -d

database-start:
	docker start postgres

database-stop:
	docker stop postgres
FROM        golang:bullseye 

ENV         GIN_MODE=release
ENV         PORT=8080
ENV         GOOS=linux
ENV         GOARCH=amd64
ENV         CGO_ENABLED=0

WORKDIR     /go/src/glass

COPY        . .

RUN         apt-get update && apt-get upgrade -y
RUN         apt-get install nano 
RUN         go get ./...
RUN         go build -o ./bin/glass ./cmd/main.go

EXPOSE      $PORT

ENTRYPOINT ["./bin/glass"]
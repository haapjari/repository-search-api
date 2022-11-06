FROM        golang:alpine

ENV         GIN_MODE=release
ENV         PORT=8080

WORKDIR     /go/src/glass

COPY        . .

RUN         apk update && apk add --no-cache git

RUN         go get ./...

RUN         go build -o ./bin/glass ./cmd/main.go

EXPOSE      $PORT

ENTRYPOINT ["./bin/glass"]
# Builder

FROM        golang:latest as build

ENV         GIN_MODE=release
ENV         GOOS=linux
ENV         GOARCH=amd64
ENV         GO111MODULE=on
ENV         CGO_ENABLED=0

WORKDIR     /workspace

COPY        . . 

RUN         apt-get update && apt-get upgrade -y  && \
            go get ./...                          && \
            go build -o glass ./cmd/main.go

# Final

FROM        scratch as final
 
COPY        --from=build /workspace/glass ./

EXPOSE      8080

ENTRYPOINT ["/glass"]

#
# Build Stage
#

FROM        golang:latest as build

ENV         GIN_MODE=release
ENV         GOOS=linux
ENV         GOARCH=amd64
ENV         GO111MODULE=on
ENV         CGO_ENABLED=0

WORKDIR     /workspace

COPY        . . 

RUN         go get ./...                          && \
            go build -o glass ./cmd/main.go

#
# Final Stage
#

FROM        scratch as final
 
COPY        --from=build /workspace/glass ./

EXPOSE      8080

ENTRYPOINT ["/glass"]

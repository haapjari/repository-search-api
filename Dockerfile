#
## Build
#

FROM golang:1.24.3 AS build

ENV GOOS=linux \
    GOARCH=amd64 \
    GO111MODULE=on \
    CGO_ENABLED=0

WORKDIR /workspace

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o /workspace/rsa ./cmd/main.go

#
## Final
#

FROM scratch

WORKDIR /workspace

COPY --from=build /workspace/rsa .

ENTRYPOINT ["./rsa"]

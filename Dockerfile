# BUILDER IMAGE

ARG     VERSION=bullseye
FROM    golang:${VERSION} as builder

WORKDIR /usr/src/

COPY    go.mod go.mod
COPY    go.sum go.sum

RUN     go mod download

COPY    . .

RUN     go build -a -o glass cmd/main.go

RUN     go install github.com/hhatto/gocloc/cmd/gocloc@latest

RUN     export GIN_MODE=release

CMD [ "./glass" ]
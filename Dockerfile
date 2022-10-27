# BUILDER IMAGE
ARG     VERSION=bullseye

FROM    golang:${VERSION} as builder

WORKDIR /usr/src/

COPY    . .

RUN     go mod download

RUN     go build -a -o glass cmd/main.go

RUN     export GIN_MODE=release

CMD [ "./glass" ]
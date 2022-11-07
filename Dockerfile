#######################################
#######################################

FROM        golang:bullseye as builder

# General Environment Variables
ENV         GIN_MODE=release
ENV         PORT=8080

# General GoLang Environment Variables
ENV         GOOS=linux
ENV         GOARCH=amd64
ENV         CGO_ENABLED=0

WORKDIR     /go/src/glass

# Copy Files from the Local Filesystem to the Image
COPY        . .

# Update Packages and Install Git
RUN         apt-get update && apt-get install git 

# Install the Dependencies
RUN         go get ./...

# Build the Binary
RUN         go build -o ./bin/glass ./cmd/main.go

#######################################
#######################################

# Scratch Image
FROM        scratch

# Environment Variables
ENV         PORT=8080

# Because it's a scratch image, other paths does not exist.
WORKDIR     /

# Copy Binary to the Scratch Image
COPY        --from=builder /go/src/glass/bin/glass /glass

# Copy Environment Variables
# TODO: This should be passed from the docker-compose
COPY        --from=builder /go/src/glass/.env /.env

# Expose the Port
EXPOSE      $PORT

# Run the Binary
ENTRYPOINT ["./glass"]
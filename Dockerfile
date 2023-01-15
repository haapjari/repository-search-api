FROM        golang:bullseye 

# General Environment Variables
ENV         GIN_MODE=release
ENV         PORT=8080
ENV         GOOS=linux
ENV         GOARCH=amd64
ENV         CGO_ENABLED=0

WORKDIR     /go/src/glass

# Copy Files from the Local Filesystem to the Image
COPY        . .

# Install vim
RUN         apt-get update && apt-get upgrade -y
RUN         apt-get install nano 

# Install the Dependencies
RUN         go get ./...

# Build the Binary
RUN         go build -o ./bin/glass ./cmd/main.go

# Expose the Port
EXPOSE      $PORT

# Run the Binary
ENTRYPOINT ["./bin/glass"]
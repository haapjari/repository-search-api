FROM        golang:bullseye 

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

# Install hhatto/gocloc Binary
RUN         go install github.com/hhatto/gocloc/cmd/gocloc@latest

# Build the Binary
RUN         go build -o ./bin/glass ./cmd/main.go

# Environment Variables
ENV         PORT=8080
ENV         GIN_MODE=release

# Expose the Port
EXPOSE      $PORT

# Run the Binary
ENTRYPOINT ["./bin/glass"]
# Start from the latest golang base image
FROM golang:1.14.1-alpine

# Add Maintainer Info
LABEL maintainer="Gustavo P Borges <gugutz@gmail.com>"

WORKDIR /api

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN go build -tags dev -o app cmd/main.go

EXPOSE 3000

CMD ["./app"]
# CMD ["go", "run", "-tags", "dev", "cmd/main.go"]

# Start from the latest golang base image
FROM golang:latest AS builder

# Add Maintainer Info
LABEL maintainer="Philipp Larionov <https://github.com/PhilLar>"

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 go build ./cmd/...

# final stage 
FROM alpine:latest



# Copy the source from the current directory to the Working Directory inside the container
COPY --from=builder /app/images-back /bin/

# Copy the source from the current directory to the Working Directory inside the container
COPY --from=builder /app/files /files/

# Copy the source from the current directory to the Working Directory inside the container
COPY --from=builder /app/migrations /migrations/

# Copy the source from the current directory to the Working Directory inside the container
COPY --from=builder /app/static /bin/static/

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
ENTRYPOINT ["/bin/images-back"]

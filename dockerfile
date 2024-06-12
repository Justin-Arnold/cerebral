# Use the official golang image as base
FROM golang:1.22-alpine AS build

# Set the curren:t working directory inside the container
WORKDIR /app

# Copy the Go modules manifests
COPY go.mod .
COPY go.sum .

# Download dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go app
RUN go build -o main .

# Start a new stage from scratch
FROM alpine:latest

# Set the current working directory inside the container
WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=build /app/main .
# Copy the templates directory into the container
COPY --from=build /app/templates ./templates
# Copy the static directory into the container
COPY --from=build /app/static ./static

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./main"]

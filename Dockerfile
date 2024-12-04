# Use the Go 1.23-alpine image
FROM golang:1.23-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files to the workspace
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go application
RUN go build -o main .

# Expose port 8080 to the host
EXPOSE 8080

# Command to run when starting the container
CMD ["./main"]

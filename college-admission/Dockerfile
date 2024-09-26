# Stage 1: Build the application
FROM golang:1.23 AS builder
WORKDIR /app

# Copy go.mod and go.sum files to cache dependencies
COPY go.mod go.sum ./
RUN go mod download
#RUN go mod tidy -v

# Copy the source code
COPY main.go .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64  go build -o college-admission-be .

# Add echo statements
RUN echo "Current directory:" && pwd
RUN echo "Files in current directory:" && ls -la .
RUN echo "Files in /app:" && ls -la /app


# Stage 2: Create a lightweight container
#FROM alpine:latest

FROM golang:1.23-alpine

WORKDIR /root/

# Add echo statements
RUN echo "Current directory:" && pwd
RUN echo "Files in current directory:" && ls -la .
RUN echo "Files in /root:" && ls -la /root

# Copy the binary from the builder stage
COPY --from=builder /app/college-admission-be .

# Expose port 8080 for the app
EXPOSE 8080

RUN chmod +x ./college-admission-be

# Run the app
CMD ["./college-admission-be"]
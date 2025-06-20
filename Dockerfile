# Start from the official Golang image for building
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
RUN cd cmd/server && go build -o /app/server

# Use a minimal image for running
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/server ./server
EXPOSE 8080
ENV HTTP_PORT=8080
CMD ["./server"] 

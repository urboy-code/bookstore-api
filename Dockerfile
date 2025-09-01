# --- Step 1: Build Stage ---
# Use the official Golang image for compiltaion
FROM golang:1.24-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy file go.mod and go.sum to the working directory for caching dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire project to the working directory
COPY . .

# Copy the Go application to be one file binary static
# CGO_ENABLED=0 and GGOS=linux important for create compatible binary
RUN CGO_ENABLED=0 GGOS=linux go build -o /bookstore-app

# --- Step 2: Run Stage ---
# Use a minimal image for running the application
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /root/

# Copy the compiled Go binary from the builder stage
COPY --from=builder /bookstore-app .

# Expose the port that the application will run on (8080)
EXPOSE 8080

# Command to run the application
CMD ["./bookstore-app"]
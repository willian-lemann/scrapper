# Use an official Go image as a builder
FROM golang:1.22.1 AS builder

# Set the working directory
WORKDIR /app

# Copy the Go module files
COPY go.mod go.sum ./

# Download Go module dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go script
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main cmd/listings/main.go

# Use a lightweight Alpine image
FROM alpine:3.6

# Install necessary packages
RUN apk --no-cache add curl bash

# Copy the Go binary
COPY --from=builder /app/main /app/main
COPY --from=builder /app/.env .  

# Ensure the Go binary has execution permissions
RUN chmod +x /app/main

# Set environment variables
ARG SUPABASE_URL
ARG SUPABASE_ANON_KEY
ARG DATABASE_URL

ENV SUPABASE_URL=${SUPABASE_URL}
ENV SUPABASE_ANON_KEY=${SUPABASE_ANON_KEY}
ENV DATABASE_URL=${DATABASE_URL}

# Run the Go binary
CMD ["/app/main"]

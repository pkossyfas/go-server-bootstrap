# Use golang 1.19 (check go.mod) as base image
FROM golang:1.19 as builder

LABEL maintainer='pkossyfas@outlook.com'

# Create the app user
RUN useradd appuser -M

WORKDIR /app

# copy go mod and sum files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# build the Go app
COPY . .
ARG VERSION=develop
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-X github.com/pkossyfas/go-server-bootstrap/controller.AppVersion=$VERSION" -a -o go-server-bootstrap -tags netgo

# Start a new stage from scratch
FROM scratch

# Copy the user from builder
COPY --from=builder /etc/passwd /etc/passwd

COPY --from=builder /app/go-server-bootstrap /usr/bin/

USER appuser

# Run the executable
CMD ["go-server-bootstrap"]

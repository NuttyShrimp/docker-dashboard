# First stage: Get Golang image from DockerHub.
FROM golang:1.26.5-alpine3.24 AS backend-builder

# Set our working directory for this stage.
WORKDIR /app

# Copy all of our files.
COPY . .

# Get and install all dependencies.
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/api/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o migrate ./migrate.go

# Last stage: discard everything except our executables.
FROM alpine:3.24 AS prod

# Set our next working directory.
WORKDIR /app

# Copy our executable and our built React application.
COPY --from=backend-builder /app/server .
COPY ./config ./config

ENV APP_ENV=production

# Declare entrypoints and activation commands.
EXPOSE 8000
ENTRYPOINT ["./server"]

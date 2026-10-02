# Backend build stage
FROM golang:1.26.5-alpine3.24 AS backend-builder

WORKDIR /app

# Cache Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy source code and compile backend binary
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/api/main.go

# Frontend / CSS build stage
FROM node:24-alpine AS frontend-builder

ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"

RUN corepack enable

WORKDIR /app

# Cache dependencies
COPY package.json pnpm-lock.yaml pnpm-workspace.yaml* ./
RUN --mount=type=cache,id=pnpm,target=/pnpm/store \
    pnpm install --frozen-lockfile

# Copy assets and views for Tailwind scanning and compilation
COPY assets ./assets
COPY views ./views

# Build CSS with Tailwind CLI via pnpm
RUN pnpm run build:css

# Final production stage
FROM alpine:3.24 AS prod

WORKDIR /app

# Copy server executable and static assets
COPY --from=backend-builder /app/server .
COPY --from=frontend-builder /app/public ./public
COPY ./config ./config

ENV APP_ENV=production

EXPOSE 3000
ENTRYPOINT ["./server"]

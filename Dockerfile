FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /build
COPY go.mod ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-w -s" \
    -o /action-semantic-versioning ./cmd/action-semantic-versioning

FROM alpine:latest

RUN apk add --no-cache git && \
    git config --global --add safe.directory /github/workspace

COPY --from=builder /action-semantic-versioning /action-semantic-versioning

ENTRYPOINT ["/action-semantic-versioning"]

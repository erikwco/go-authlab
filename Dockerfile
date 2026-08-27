# Multi-stage build: compile a static binary, then run as non-root.
FROM golang:1.22-bookworm AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/authlab ./cmd/authlab

FROM alpine:3.20

RUN apk add --no-cache ca-certificates \
    && adduser -D -H -u 65532 authlab

COPY --from=builder /out/authlab /usr/local/bin/authlab

USER authlab
EXPOSE 8080

ENTRYPOINT ["/usr/local/bin/authlab"]

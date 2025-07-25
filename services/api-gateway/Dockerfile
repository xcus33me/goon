# Modules caching
FROM golang:1.24.2-alpine3.21 as modules

COPY go.mod go.sum /modules/

WORKDIR /modules

RUN go mod download

# Builder
FROM golang:1.24.2-alpine3.21 as builder

COPY --from=modules /go/pkg go/pkg

COPY . /app

WORKDIR /app

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o /bin/app ./cmd/app

FROM alpine:3.21

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /app/config /config
COPY --from=builder /bin/app /app/app

WORKDIR /app

CMD ["/app/app"]

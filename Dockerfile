ARG GO_VERSION=1.25
FROM golang:${GO_VERSION}-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG SERVICE
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /build/service ./services/${SERVICE}

FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /build/service /service

EXPOSE 8080

ENTRYPOINT ["/service"]

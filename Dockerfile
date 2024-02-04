FROM golang:1.19 AS builder

WORKDIR /src

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o ./server ./cmd

FROM alpine:3

ENV ENDPOINT_CONFIG_PATH=/app/config.yaml

WORKDIR /app

COPY --from=builder /src/server /src/config/config.yaml ./

EXPOSE 8080

ENTRYPOINT ["/app/server"]

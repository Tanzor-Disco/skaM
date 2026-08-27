FROM golang:latest AS builder

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o backend .

FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/backend .

RUN apt-get update && apt-get install -y ca-certificates

CMD ["./backend"]

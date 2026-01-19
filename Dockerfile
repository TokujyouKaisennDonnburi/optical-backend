FROM golang:1.24-bookworm AS builder

WORKDIR /app

RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates

COPY go.mod go.sum ./

RUN go mod download

COPY . .

ENV CGO_ENABLED=0

RUN CGO_ENABLED=0 GOOS=linux go build -buildvcs=false -o /usr/local/bin/go_app ./cmd/api/

FROM scratch
WORKDIR /app

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/local/bin/go_app /app/go_app
COPY --from=builder /app/db /app/db

EXPOSE 8000
ENTRYPOINT ["./go_app"]

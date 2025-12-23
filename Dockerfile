FROM golang:1.24-bookworm AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

ENV CGO_ENABLED=0

RUN CGO_ENABLED=0 GOOS=linux go build -buildvcs=false -o /usr/local/bin/go_app ./cmd/api/

FROM scratch
WORKDIR /app

COPY --from=builder /usr/local/bin/go_app /app/go_app
COPY --from=builder /app/db /app/db

ENTRYPOINT ["./go_app"]

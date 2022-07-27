# Build stage
FROM golang:1.18-alpine3.16 AS builder
WORKDIR /app
RUN apk add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz | tar xvz
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
# RUN CGO_ENABLED=0 go test -v ./...
RUN go build -o main main.go

# Run stage
FROM alpine:3.16
WORKDIR /app
COPY --from=builder /app/migrate ./migrate
COPY --from=builder /app/main .
COPY app.env .
COPY db/migration ./migration
COPY start.sh .

EXPOSE 8080
CMD ["/app/main"]
ENTRYPOINT [ "/app/start.sh" ]

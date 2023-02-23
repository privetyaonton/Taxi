FROM golang:alpine AS builder

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o ./bin/main ./cmd/main.go

FROM scratch

WORKDIR /app

COPY --from=builder /app/bin/main .

COPY . .

EXPOSE 8080

CMD ["./main"]
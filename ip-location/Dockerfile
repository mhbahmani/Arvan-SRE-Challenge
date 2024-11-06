FROM docker.arvancloud.ir/golang:1.22.5-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o ip-location .

FROM docker.arvancloud.ir/alpine:3.18

WORKDIR /app

COPY --from=builder /app/ip-location .

EXPOSE 8080

CMD ["./ip-location"]

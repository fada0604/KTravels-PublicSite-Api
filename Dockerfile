FROM golang:1.26-alpine AS builder

RUN apk add --no-cache curl ca-certificates

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download || (go mod tidy && go mod download)

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o ./api ./app/main.go

FROM alpine:3.19

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /app/api .

EXPOSE 8080

ENV APP_ENV=production

CMD ["/app/api"]
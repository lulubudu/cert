FROM golang:latest as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY ./ ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /notifier .

FROM alpine:latest

COPY --from=builder /notifier /

CMD ["/notifier"]
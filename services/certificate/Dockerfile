FROM golang:latest as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY ./ ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /certificate .


FROM alpine:latest

COPY --from=builder /certificate /

EXPOSE 8080

CMD ["/certificate"]
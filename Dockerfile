FROM golang:1.24-alpine

WORKDIR /app

COPY . .

RUN  go build -o bloGo ./cmd/web/*.go

CMD ["./go-chatty"]
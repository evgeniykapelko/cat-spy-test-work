# Dockerfile
FROM golang:1.22.5-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main ./cmd/app

EXPOSE 8087

CMD ["./main"]

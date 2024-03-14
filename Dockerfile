FROM golang:latest AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

WORKDIR /app
COPY . .
RUN chmod +x run.sh

EXPOSE 8080

CMD ["./run.sh"]
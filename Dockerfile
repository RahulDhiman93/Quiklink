FROM golang:latest AS builder

WORKDIR /app
COPY go.mod go.sum ./

RUN go mod download

COPY . .

WORKDIR /app

COPY . .

RUN chmod +x prod.sh

EXPOSE 8080

CMD ["sh", "prod.sh"]
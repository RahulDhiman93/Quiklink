FROM golang:latest AS builder
WORKDIR /var/www/app
COPY go.mod go.sum ./
RUN go mod download
WORKDIR /var/www/app
COPY . .
RUN chmod +x run.sh
EXPOSE 8080
CMD ["sh", "run.sh"]


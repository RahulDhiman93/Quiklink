FROM postgres AS postgres_setup
COPY init_db.sql /docker-entrypoint-initdb.d/
FROM postgres_setup AS final

FROM golang:latest AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go install github.com/gobuffalo/pop/v6/soda@latest

WORKDIR /app
COPY . .
RUN chmod +x run.sh

EXPOSE 8080

CMD ["./run.sh"]
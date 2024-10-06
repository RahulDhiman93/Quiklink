FROM golang:alpine3.19 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

WORKDIR /app
COPY . .

RUN go build -o Quiklink_App cmd/*.go
RUN chmod +x ./Quiklink_App

RUN apk update && apk add --no-cache \
    dialog \
    openssh-server \
    caddy

RUN echo "root:Docker!" | chpasswd \
    && chmod u+x ./entrypoint.sh

COPY sshd_config /etc/ssh/

# Configure Caddy
RUN echo ":8080 {\n  reverse_proxy localhost:8080\n}" > /etc/caddy/Caddyfile

EXPOSE 8080 2222

CMD ["sh", "-c", "caddy run --config /etc/caddy/Caddyfile & ./Quiklink_App"]
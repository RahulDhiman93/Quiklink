FROM golang:alpine3.19 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
WORKDIR /app
COPY . .
RUN go build -o Quiklink_App cmd/*.go
RUN chmod +x ./Quiklink_App
EXPOSE 8080
CMD ["./Quiklink_App"]
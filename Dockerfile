FROM golang:latest AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o Quiklink_BE cmd/*.go

RUN go install github.com/gobuffalo/pop/v6/soda@latest

FROM ubuntu:22.04
RUN echo 'APT::Install-Suggests "0";' >> /etc/apt/apt.conf.d/00-docker
RUN echo 'APT::Install-Recommends "0";' >> /etc/apt/apt.conf.d/00-docker
ENV TZ=America/New_York
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone
RUN DEBIAN_FRONTEND=noninteractive \
  apt-get update \
  && apt-get install -y postgresql-14 \
  && rm -rf /var/lib/apt/lists/*

RUN service postgresql start

WORKDIR /app
COPY --from=builder /app/Quiklink_BE .
COPY --from=builder /go/bin/soda /usr/local/bin/soda
COPY migrations/ /app/migrations/

EXPOSE 8080

RUN soda migrate up

CMD ["./Quiklink_BE", "-dbname=Quiklink_BE", "-dbuser=rahuldhiman"]
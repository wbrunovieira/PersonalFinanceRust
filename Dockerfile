FROM golang:1.20


WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY db /app/db


RUN apt-get update && \
  apt-get install -y --no-install-recommends \
  ca-certificates \
  curl \
  && curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz \
  | tar xvz && mv migrate /usr/local/bin/ && \
  chmod +x /usr/local/bin/migrate && \
  apt-get clean && rm -rf /var/lib/apt/lists/*

RUN go build -o /app/main .

EXPOSE 8080

CMD ["/app/main"]

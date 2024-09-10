FROM golang:latest

WORKDIR /

COPY go.mod .

RUN go mod download

COPY . .

RUN go build -o service ./cmd/cli

EXPOSE 8000

CMD ["./service"]

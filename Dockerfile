FROM golang:1.20

RUN apt-get update && apt-get install -y postgresql-client

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o ./build/app ./cmd/app.go

CMD ["./build/app"]

FROM golang:1.25

RUN apt-get update

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o ./build/app ./cmd/app.go

CMD ["./build/app"]

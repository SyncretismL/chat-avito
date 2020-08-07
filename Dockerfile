FROM golang:1.14-alpine

WORKDIR /app
COPY go.mod go.sum logs.log ./

RUN go mod download

COPY . .

RUN go build -o /main ./cmd

EXPOSE 9000

CMD ["/main"]
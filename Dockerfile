
FROM golang:1.21.3

RUN mkdir /app

WORKIRD /app

RUN go mod tidy

RUM go build main.go

CMD ["./main"]
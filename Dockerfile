FROM golang:1.21.3

RUN mkdir /app

COPY ./ /app

WORKDIR /app

RUN go mod tidy

RUN go build ./main.go

RUN CGO_ENABLED=0 GOOS=linux go build -o /main

#FROM scratch

EXPOSE 8081

CMD ["/main"]


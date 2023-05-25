FROM golang:alpine

WORKDIR /testtask
COPY . .

RUN go build -ldflags '-w -s' -a -o testtask main.go

CMD ["./testtask"]
EXPOSE 8080
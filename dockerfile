FROM golang:latest

COPY ./ ./

RUN mkdir ./config
RUN go build ./cmd/main.go

CMD ["./main"]


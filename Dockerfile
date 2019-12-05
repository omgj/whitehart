FROM golang:latest
WORKDIR /Users/oliver/go/whitehart
COPY . .

RUN go mod download
RUN  go build -o main .

EXPOSE 8080

CMD ["./main"]

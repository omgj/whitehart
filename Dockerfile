FROM golang:latest

WORKDIR public.html /whitehart

COPY . .

RUN go build -o main .

EXPOSE 8080

CMD ["./main"]

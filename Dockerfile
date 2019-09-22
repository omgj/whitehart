FROM golang:latest
WORKDIR /Users/oliver/go/whitehart
COPY . .


RUN go get cloud.google.com/go && go get cloud.google.com/go/firestore && go build -o main .

EXPOSE 8080

CMD ["./main"]
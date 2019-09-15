FROM golang:latest

WORKDIR cloudbuild.yaml public.html /whitehart

COPY . .

RUN go build -o main .

EXPOSE 8080

CMD ["./main"]

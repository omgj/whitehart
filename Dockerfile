FROM golang:1.13-buster as builder
ENV GO111MODULE=on
WORKDIR /app
COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /go/bin/app

FROM gcr.io/distroless/base-debian10
COPY --from=build /go/bin/app /go/bin/app

EXPOSE 8080

CMD ["app"]

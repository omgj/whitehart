FROM golang:1.13-buster as builder
ENV GO111MODULE=on
WORKDIR /Users/oliver/go/whitehart
COPY go.mod /go/src/app
COPY go.sum /go/src/app
RUN go mod download

COPY . /go/src/app

RUN go build

FROM gcr.io/distroless/base-debian10
COPY --from=build /go/bin/app /

EXPOSE 8080

CMD ["app"]

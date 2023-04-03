FROM golang:alpine as builder
RUN apk add git

WORKDIR /go/src/app

ENV GO111MODULE=auto

RUN go get github.com/cespare/reflex

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o ./run .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

# Copy the executable from the builder
COPY --from=builder /go/src/app/run .

CMD ["./run"]
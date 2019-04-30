FROM golang:alpine

RUN apk add --update --no-cache alpine-sdk bash ca-certificates \
      libressl \
      tar \
      git openssh openssl yajl-dev zlib-dev cyrus-sasl-dev openssl-dev build-base coreutils

WORKDIR $GOPATH/src/github.com/tarosaiba/kafka-train-producer

COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["go", "run", "server.go"]

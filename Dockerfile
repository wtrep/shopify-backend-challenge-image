FROM golang:1.15.1

ENV GOPRIVATE=github.com/wtrep/*

WORKDIR /go/src/app

COPY image ./image
COPY common ./common
COPY main.go .
COPY go.mod .
COPY go.sum .

RUN go get -v -d
RUN go install -v

EXPOSE 8080
VOLUME /opt/certs

CMD ["shopify-backend-challenge-image"]
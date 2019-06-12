FROM golang:1.11

ENV GOOS=linux

WORKDIR $GOPATH/src/broker

COPY . .

ENV GO111MODULE=on

RUN go mod download
RUN go build
RUN chmod +x ./dcm-middleware

CMD ["./dcm-middleware"]
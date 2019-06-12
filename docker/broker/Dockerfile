FROM golang:1.11

LABEL maintainer="Luiz Reis <lrnn@cin.ufpe.br>"

ENV GOOS=linux

WORKDIR $GOPATH/src/broker

COPY . .

ENV GO111MODULE=on

RUN go mod download
RUN go build
RUN chmod +x ./projetomiddleware

CMD ["./projetomiddleware"]
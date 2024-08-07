FROM golang:latest


RUN mkdir /golang

RUN go install  github.com/air-verse/air@latest

ADD . /golang/

RUN go install  github.com/air-verse/air@latest

WORKDIR /golang

RUN go mod download

CMD ["air", "-c", ".air.toml"]
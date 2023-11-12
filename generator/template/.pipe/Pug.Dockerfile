FROM golang:latest

WORKDIR /vectra

RUN go install github.com/Joker/jade/cmd/jade@latest

CMD [ "sleep", "infinity" ]

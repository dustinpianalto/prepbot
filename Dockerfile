FROM golang:1.15-alpine as dev

WORKDIR /go/src/prepbot
COPY ./go.mod .
COPY ./go.sum .

RUN go mod download

COPY . .
RUN go install github.com/dustinpianalto/prepbot/...

CMD [ "go", "run", "cmd/prepbot/main.go"]

from alpine

WORKDIR /bin

COPY --from=dev /go/bin/prepbot ./prepbot

CMD [ "prepbot" ]

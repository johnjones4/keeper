FROM golang:1.21

WORKDIR /usr/src/app
COPY . .

RUN go get .
RUN go build --tags "fts5" .

CMD ["/usr/src/app/main"]

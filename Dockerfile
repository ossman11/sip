FROM golang:1.8

WORKDIR /go/src/github.com/ossman11/sip
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

EXPOSE 1670

CMD [ "sip" ]

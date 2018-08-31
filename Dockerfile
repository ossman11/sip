FROM golang:1.8

ENV PORT=1670

WORKDIR /go/src/github.com/ossman11/sip
COPY . .

RUN ./crt/make
RUN go get -d -v ./...
RUN go install -v ./...

EXPOSE 1670

CMD [ "sip" ]

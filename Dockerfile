FROM golang:alpine

WORKDIR /go/src/feedr
COPY . .
RUN go get -d -v ./...
RUN go install -v ./...

ENTRYPOINT ["feedr"]

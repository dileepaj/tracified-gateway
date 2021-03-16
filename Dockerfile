FROM golang:1.15.6
LABEL maintainer = "Nisal Perera <nisaledu@gmail.com>"
RUN mkdir -p /go/src/github.com/dileepaj/tracified-gateway/
COPY . /go/src/github.com/dileepaj/tracified-gateway/
WORKDIR /go/src/github.com/dileepaj/tracified-gateway/
RUN go get -u github.com/golang/dep/cmd/dep
#RUN dep init
RUN dep ensure
RUN go build
CMD ["./tracified-gateway"]

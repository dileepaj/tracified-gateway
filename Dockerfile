FROM golang:1.15.6
LABEL maintainer = "Nisal Perera <nisaledu@gmail.com>"
RUN mkdir -p /go/src/github.com/dileepaj/tracified-gateway/
COPY . /go/src/github.com/dileepaj/tracified-gateway/
WORKDIR /go/src/github.com/dileepaj/tracified-gateway/
RUN go get -u github.com/golang/dep/cmd/dep
#RUN dep init
ENV PublicKey=""
ENV SecretKey=""
ENV GATEWAY_PORT=""
ENV BRANCH_NAME=""
ENV DBUSERNAME=""
ENV DBPASSWORD=""
ENV DBHOST=""
ENV DBPORT=""
ENV DBNAME=""
ENV ADMINDBUSERNAME=""
ENV ADMINDBPASSWORD=""
ENV ADMINDBHOST=""
ENV ADMINDBPORT=""
ENV ADMINDBNAME=""
RUN dep ensure
RUN go build
CMD ["./tracified-gateway"]

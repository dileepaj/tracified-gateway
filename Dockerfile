FROM golang:1.18-alpine
RUN apk add --update cmake gcc g++ git  make  tar wget python3

# Set destination for COPY
RUN mkdir -p /go/src/github.com/dileepaj/tracified-gateway/
WORKDIR /go/src/github.com/dileepaj/tracified-gateway/

# Download Go modules
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the source code
COPY . ./

# Build
RUN go build github.com/dileepaj/tracified-gateway
COPY . ./

# CMD [ "ls" ]
# RUN ls
RUN chmod +x tracified-gateway
# RUN find . -type f | grep "tracified-gateway"
CMD ["./tracified-gateway"]

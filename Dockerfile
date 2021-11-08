FROM golang:1.16-alpine

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
RUN CGO_ENABLED=0 go build github.com/dileepaj/tracified-gateway
COPY . ./

# CMD [ "ls" ]
# RUN ls
RUN chmod +x tracified-gateway
# RUN find . -type f | grep "tracified-gateway"
CMD ["./tracified-gateway"]

FROM alpine:3.9

# Installs
RUN apk add --no-cache musl-dev git go

# Configure Go
ENV GOROOT /usr/lib/go
ENV GOPATH /go
ENV PATH /go/bin:$PATH
RUN mkdir -p ${GOPATH}/src ${GOPATH}/bin

# Verify
RUN go version

# Install Vision API
RUN go get -u cloud.google.com/go/vision/apiv1

# Get source code
WORKDIR $GOPATH
RUN git clone https://github.com/gabrielpsilva/GoogleVisionAPI-Go.git

# Set image entry point
ENTRYPOINT go run $GOPATH/GoogleVisionAPI-Go/main.go

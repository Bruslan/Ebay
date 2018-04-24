FROM ubuntu:16.04
FROM golang

# Copy the local package files to the container's workspace.
ADD . /go/src/malte
ENV GOBIN=/go/bin
# WORKDIR /go/src

# Build the outyet command inside the container.
# RUN go get....
RUN go get golang.org/x/net/http2
# RUN git clone https://github.com/jplaui/malte.git
RUN go install /go/src/malte/server.go

# Document that the service listens on port 443.
EXPOSE 443
EXPOSE 80

WORKDIR /go/bin
CMD ["./server"]


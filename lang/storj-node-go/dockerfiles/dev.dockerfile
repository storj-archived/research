FROM golang

# install build dependencies
RUN go get -u github.com/golang/dep/cmd/dep

# copy source
RUN mkdir -p /go/src/storj-node-go
COPY ./Gopkg.lock ./Gopkg.toml ./main.go /go/src/storj-node-go/

WORKDIR /go/src/storj-node-go

# install go dependencies
RUN dep ensure

# volume mount src code
VOLUME /go/src/storj-node-go

# expose port for server
EXPOSE 8080

#CMD go run main.go
CMD ["/bin/bash"]

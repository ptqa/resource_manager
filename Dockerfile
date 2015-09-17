FROM golang:onbuild
WORKDIR /go/src/github.com/ptqa/resource_manager
ADD . /go/src/github.com/ptqa/resource_manager
ADD run_docker.sh /tmp/
RUN go get github.com/ptqa/resource_manager
CMD ["/tmp/run_docker.sh"]

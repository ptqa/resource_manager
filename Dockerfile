FROM golang:onbuild
WORKDIR /go/src/github.com/ptqa/resource_manager
ADD . /go/src/github.com/ptqa/resource_manager
ADD run_docker.sh /tmp/
RUN go get github.com/ptqa/resource_manager
EXPOSE 3000
CMD ["/tmp/run_docker.sh"]

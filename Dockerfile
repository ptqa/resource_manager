FROM golang:onbuild
WORKDIR /gopath/src/github.com/ptqa/resource_manager
ADD . /gopath/src/github.com/ptqa/resource_manager
RUN go get github.com/ptqa/resource_manager
#ENTRYPOINT ["/gopath/bin/app"]
#CMD ["/gopath/bin/resource_manager"]

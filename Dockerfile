#source
FROM golang:latest
#author
MAINTAINER luopeng "755200@qq.com"
#workdir
WORKDIR $GOPATH/
#add code
ADD hello.go $GOPATH/
#build
RUN go build hello.go
#expose
EXPOSE 8081
#entrypoint
ENTRYPOINT  ["./hello"]

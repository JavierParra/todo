FROM golang
RUN mkdir -p /go/src
ADD src/ /go/src
RUN go get github.com/nu7hatch/gouuid
WORKDIR /app
RUN go build -o main /go/src/server.go

CMD ["/app/main"]

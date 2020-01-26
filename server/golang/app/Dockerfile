FROM golang
RUN go get github.com/nu7hatch/gouuid
RUN go get github.com/githubnemo/CompileDaemon

ENV ROOT /app
ENV SRCPATH ${ROOT}/src/server
ENV GOPATH "${GOPATH}:${ROOT}"

RUN echo $GOPATH
RUN mkdir -p ${SRCPATH}
ADD src ${SRCPATH}
WORKDIR ${ROOT}
# RUN go build -o server /go/src/server.go

# CMD ["/app/server"]
CMD ["sh",  "-cx", \
	"CompileDaemon -verbose -directory=$SRCPATH -build=\"go build -o ${ROOT} \" -color -command=./server" \
]

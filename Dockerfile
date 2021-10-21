FROM golang:alpine
WORKDIR /app
COPY ./go.mod .
RUN go env -w GOPROXY=goproxy.cn,direct && go mod download
COPY . .
CMD ["go", "build", "-o", "event_retriever" , "src/server.go"]
CMD ["./event_retriever"]
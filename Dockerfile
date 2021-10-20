FROM golang:alpine
WORKDIR /app
COPY ./go.mod .
RUN go mod download
RUN go build -o event_retriever src/server.go
COPY . .
CMD ["./event_retriever"]
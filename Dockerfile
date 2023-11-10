#FIRST STAGE

FROM golang:latest as builder
WORKDIR /app
ENV GOPROXY https://goproxy.io 
#use the goproxy.io proxy server for downloading Go module dependencies
COPY go.mod go.sum ./
#o.mod and go.sum files are copied to the working directory.
RUN go mod download
# go mod download command is executed to download the Go module dependencies.
COPY . .
#entire application code is copied to the working directory.
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o livego . 
#disable CGO (C Go) support. Operating system as Linux.
#build the Go application. The resulting binary is named livego and is saved in the current directory.


#SECOND STAGE
FROM alpine:latest
RUN mkdir -p /app/config
WORKDIR /app
ENV RTMP_PORT 1935
ENV HTTP_OPERATION_PORT 8090
COPY --from=builder /app/livego .
EXPOSE ${RTMP_PORT}
EXPOSE ${HTTP_OPERATION_PORT}
ENTRYPOINT ["./livego"]

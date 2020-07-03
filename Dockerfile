# start a golang base image, version 1.8
FROM golang:alpine AS build

#switch to our app directory
RUN mkdir -p /build
WORKDIR /build
#copy the source files
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY main.go .
ENV CGO_ENABLED=0
#build the binary with debug information removed
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -a -o app main.go
RUN apk update && apk add upx
RUN upx /build/app


FROM scratch
ENV GIN_MODE=release
ENV PORT=8000
COPY --from=build /build/app /app
COPY assets /assets
COPY static /static

WORKDIR /

ENTRYPOINT ["/app"]

FROM golang:1.19.4-alpine AS build_base

ENV CGO_ENABLED=1
ENV GO111MODULE=on
RUN apk add --no-cache git gcc g++

# Set the Current Working Directory inside the container
WORKDIR /src

# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

# Build the Go app
RUN go build -o ./out/app ./cmd/main.go

# Start fresh from a smaller image
FROM alpine:3.17.0
RUN apk add ca-certificates

WORKDIR /app

COPY --from=build_base /src/out/app /app/webapi
COPY --from=build_base /src/uploadfile /app/uploadfile
COPY --from=build_base /src/data /app/data
COPY --from=build_base /src/log /app/log

RUN chmod +x webapi

# This container exposes port 8080 to the outside world
EXPOSE 8080

# Run the binary program produced by `go install`
ENTRYPOINT ./webapi
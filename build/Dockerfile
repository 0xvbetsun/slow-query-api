# Build's image
FROM golang:alpine AS build

RUN mkdir /app
WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN GO111MODULE=on CGO_ENABLED=0 go build -a -o slow-query-api

## Runtime's Image
FROM alpine

COPY --from=build /app/slow-query-api .

ENTRYPOINT ["./slow-query-api"]

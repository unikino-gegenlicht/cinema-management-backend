FROM golang:alpine AS build-service
COPY . /tmp/src
WORKDIR /tmp/src
RUN mkdir -p /tmp/build
RUN go mod download
RUN go build -o /tmp/build/app

FROM alpine:latest
COPY --from=build-service /tmp/build/app /backend
ENTRYPOINT ["/backend"]
EXPOSE 8000

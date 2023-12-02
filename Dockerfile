#build stage
FROM golang:alpine AS builder
ARG VERSION
RUN apk add --no-cache git
WORKDIR /go/src/app
COPY . .
RUN go get -d -v ./...
RUN go build -o /go/bin/app -v ./...

#final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /go/bin/app /app
COPY --from=builder /go/src/app/VERSION /VERSION
ENV BACKEND=$BACKEND
ENV HOSTNAME=$HOSTNAME
ENV TEMPLATEDIRECTORY=$TEMPLATEDIRECTORY
ENV TEMPLATETYPE=$TEMPLATETYPE
ENV REFRESHCACHE=$REFRESHCACHE
ENV PROFILER=$PROFILER
ENTRYPOINT ["/app"]
LABEL Name=zehd Version="$(cat VERSION)"
EXPOSE 80


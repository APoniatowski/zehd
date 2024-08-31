# build stage
FROM golang:alpine AS builder
ARG VERSION
WORKDIR /go/src/app
COPY . .
RUN mkdir -p /go/bin/ && \
  go mod tidy && \
  go get -d -v ./... && \
  go build -o /go/bin/app -v ./cmd/zehd/main.go

# final stage
FROM scratch
COPY --from=builder /go/bin/app /app
COPY --from=builder /go/src/app/VERSION /VERSION

# Container config
ENV BACKEND=${BACKEND} \
  NEWHOSTNAME=${NEWHOSTNAME} \
  # Templating config
  TEMPLATEDIRECTORY=${TEMPLATEDIRECTORY} \
  TEMPLATETYPE=${TEMPLATETYPE} \
  REFRESHCACHE=${REFRESHCACHE} \
  # Logging config
  APP_NAME=${APP_NAME} \
  LOGLOCATION=${LOGLOCATION} \
  LOGLEVEL=${LOGLEVEL} \
  # Git config
  GITLINK=${GITLINK} \
  GITTOKEN=${GITTOKEN} \
  GITUSERNAME=${GITUSERNAME} \
  # Paths config
  JSPATH=${JSPATH} \
  CSSPATH=${CSSPATH} \
  DOWNLOADSPATH=${DOWNLOADSPATH} \
  IMAGESPATH=${IMAGESPATH} \
  PROFILER=${PROFILER}

ARG VERSION
LABEL org.opencontainers.image.name="zehd" \
  org.opencontainers.image.version="${VERSION}"

EXPOSE 80

ENTRYPOINT ["/app"]

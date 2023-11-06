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
ENTRYPOINT ["/app"]
LABEL Name=zehd Version="$(cat VERSION)"
EXPOSE 80

# syntax=docker/dockerfile:1

################################################################################################################################
# FROM golang:alpine AS builder
# LABEL MAINTAINER="Adam Poniatowski <adaml.poniatowski@gmail.com>"
# LABEL MICROSERVICE="frontend"


# # # Download and install the latest release of dep
# ADD https://github.com/golang/dep/releases/download/v0.5.4/dep-linux-amd64 /usr/bin/dep
# RUN chmod +x /usr/bin/dep

# # install git, openssh-client and configure (if it resides in a private repo)
# # openssh-client and the git config line is optional
# RUN apk add git openssh-client
# RUN git config --global url."ssh://git@git.poniatowski.dev/adam/*".insteadOf https://git.poniatowski.dev/adam/*

# # As this is a private repo, copy a key that has been generated and public key added to github
# RUN mkdir /root/.ssh/
# ADD ~/.ssh/id_rsa /root/.ssh/id_rsa
# RUN chmod 600 /root/.ssh/id_rsa
# RUN ssh-keyscan git.poniatowski.dev > /root/.ssh/known_hosts

# # Copy the code from the host and compile it
# WORKDIR $GOPATH/src/git.poniatowski.dev/adam/zehd-frontend
# RUN git clone git@git.poniatowski.dev:adam/zehd-frontend.git .
# RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo -o /frontend-bin .
# RUN tar -zcf /package.tar.gz ./templates

# # # Copy the binary to new container and run
# FROM alpine
# RUN mkdir -p ~/frontend/templates
# WORKDIR ~
# COPY --from=builder /package.tar.gz .
# COPY --from=builder /frontend-bin .
# RUN mkdir -p /var/frontend/templates
# RUN tar -zxf package.tar.gz -C /var/frontend
# RUN rm package.tar.gz
# RUN setcap 'cap_net_bind_service,cap_sys_time,cap_sys_resource,cap_sys_nice+ep' ./frontend-bin
# ENTRYPOINT ["./frontend-bin"]

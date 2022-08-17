FROM golang:alpine as builder

MAINTAINER haimait

WORKDIR /work

RUN go env -w GOPROXY=https://goproxy.cn,direct && go env -w CGO_ENABLED=0
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o go-admin main.go


FROM alpine:latest
# change timezone
RUN apk add --no-cache tzdata && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && echo "Asia/Shanghai" >  /etc/timezone

WORKDIR /go-admin-api

COPY --from=builder /work/go-admin ./
COPY --from=builder /work/config ./config
COPY --from=builder /work/static ./static
COPY --from=builder /work/temp ./temp


EXPOSE 8000

CMD ["./go-admin","server","-c", "config/settings.dev.yml"]

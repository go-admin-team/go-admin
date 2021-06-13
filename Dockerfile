FROM alpine

ENV GOPROXY https://goproxy.cn/

RUN apk update --no-cache
RUN apk add --no-cache ca-certificates
RUN apk add --no-cache tzdata
ENV TZ Asia/Shanghai

COPY ./main /main
COPY ./config/settings.prod.yml /config/settings.yml
EXPOSE 8000
RUN  chmod +x /main
CMD ["/main","server","-c", "/config/settings.yml"]
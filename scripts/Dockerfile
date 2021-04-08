FROM alpine

COPY ./go-admin /
EXPOSE 8000

CMD ["/go-admin","server","-c", "/config/settings.yml"]

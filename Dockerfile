FROM alpine

# ENV GOPROXY https://goproxy.cn/

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories

# Runtime packages only.
#
# gcc and g++ used to be installed here and were 273MB of a 381MB image. The
# binary this image runs is compiled and linked before the image is built and
# arrives as a COPY, so nothing in the container ever invokes a compiler -
# there is no toolchain to drive it with either, since Go itself is not
# installed.
#
# That layer was also why a host could not share storage between images. apk
# resolves against an index that moves, so the layer digest differed on every
# build and no two images shared it: a host keeping one image per deployed
# commit stored a private 273MB copy each time. 68 of them filled the disk
# and the next deployment could not pull.
#
# libc6-compat stays. Nothing measured needs it - the binary CI produces is
# statically linked, and a container built without libc6-compat resolves a
# hostname and opens a database connection exactly as one built with it - but
# it is half a megabyte and it covers a ./main that was linked dynamically,
# which this Dockerfile has no way to check.
RUN apk add --no-cache ca-certificates tzdata libc6-compat
ENV TZ Asia/Shanghai

COPY ./main /main
COPY ./config/settings.demo.yml /config/settings.yml
COPY ./go-admin-db.db /go-admin-db.db
EXPOSE 8000
RUN  chmod +x /main
CMD ["/main","server","-c", "/config/settings.yml"]

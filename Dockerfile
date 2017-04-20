FROM scratch
MAINTAINER Yves Brissaud <yves.brissaud@gmail.com>
COPY plop /app
ENTRYPOINT ["/app"]

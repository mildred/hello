FROM scratch
MAINTAINER Yves Brissaud <yves.brissaud@gmail.com>
COPY app /app
ENTRYPOINT ["/app"]

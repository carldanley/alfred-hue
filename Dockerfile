FROM alpine as alpine

RUN apk --no-cache add tzdata zip ca-certificates

WORKDIR /usr/share/zoneinfo
RUN zip -r -0 /zoneinfo.zip .
ENV ZONEINFO /zoneinfo.zip

WORKDIR /
ADD alfred-hue /bin/

ENTRYPOINT [ "/bin/alfred-hue" ]

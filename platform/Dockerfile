FROM alpine
RUN apk update --no-cache && apk add --no-cache tzdata
ENV TZ Asia/Shanghai
WORKDIR /opt
COPY ./bin/cloud /opt/platform
COPY ./configs/config.yaml /opt/configs.yaml
EXPOSE 8010

CMD ./platform --conf configs.yaml
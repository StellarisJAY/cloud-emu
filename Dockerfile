FROM golang:1.22-bookworm AS builder
WORKDIR /opt/cloudemu
COPY ./api /opt/cloudemu/api
COPY ./gamesrv /opt/cloudemu/gamesrv
COPY ./third_party /opt/cloudemu/third_party
COPY go.mod /opt/cloudemu/go.mod
COPY ./common /opt/cloudemu/common
COPY ./emulator /opt/cloudemu/emulator
COPY ./platform /opt/cloudemu/platform

ENV GOPROXY https://goproxy.cn,direct
RUN rm -r /etc/apt/sources.list.d
RUN echo "deb http://mirrors.tuna.tsinghua.edu.cn/debian/ bookworm main contrib non-free non-free-firmware" >> /etc/apt/sources.list
RUN echo "deb http://mirrors.tuna.tsinghua.edu.cn/debian/ bookworm-updates main contrib non-free non-free-firmware" >> /etc/apt/sources.list
RUN echo "deb http://mirrors.tuna.tsinghua.edu.cn/debian/ bookworm-backports main contrib non-free non-free-firmware" >> /etc/apt/sources.list
RUN echo "deb http://mirrors.tuna.tsinghua.edu.cn/debian-security bookworm-security main contrib" >> /etc/apt/sources.list
RUN apt update
RUN apt install -y libx264-dev libvpx-dev libopusfile-dev libogg-dev
RUN go mod tidy
RUN cd /opt/cloudemu/gamesrv && make build
RUN cd /opt/cloudemu/platform && make build

# gamesrv服务依赖部分C语言库，需要和编译环境版本保持一致
FROM debian:bookworm-slim AS gamesrv
RUN rm -r /etc/apt/sources.list.d
RUN echo "deb http://mirrors.tuna.tsinghua.edu.cn/debian/ bookworm main contrib non-free non-free-firmware" >> /etc/apt/sources.list
RUN echo "deb http://mirrors.tuna.tsinghua.edu.cn/debian/ bookworm-updates main contrib non-free non-free-firmware" >> /etc/apt/sources.list
RUN echo "deb http://mirrors.tuna.tsinghua.edu.cn/debian/ bookworm-backports main contrib non-free non-free-firmware" >> /etc/apt/sources.list
RUN echo "deb http://mirrors.tuna.tsinghua.edu.cn/debian-security bookworm-security main contrib" >> /etc/apt/sources.list
RUN apt update
RUN apt install -y libx264-dev libvpx-dev libopusfile-dev libogg-dev
ENV TZ Asia/Shanghai
WORKDIR /opt/cloudemu
COPY --from=builder /opt/cloudemu/gamesrv/bin/cloud /opt/cloudemu/gamesrv
COPY ./configs/gamesrv.yaml /opt/cloudemu/configs.yaml
ARG PORT=9020
EXPOSE $PORT
ARG UDP_PORT_MIN=60000
ARG UDP_PORT_MAX=65535
# gamesrv服务webrtc使用的udp端口范围
EXPOSE $UDP_PORT_MIN-$UDP_PORT_MAX/udp
CMD ./gamesrv --conf configs.yaml

FROM alpine AS platform
RUN apk update --no-cache && apk add --no-cache tzdata
ENV TZ Asia/Shanghai
WORKDIR /opt/cloudemu
COPY --from=builder /opt/cloudemu/platform/bin/cloud /opt/cloudemu/platform
COPY ./configs/platform.yaml /opt/cloudemu/configs.yaml
ARG PORT=8010
EXPOSE $PORT
CMD ./platform --conf configs.yaml
FROM golang:latest AS build

# 配置模块代理
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn,direct

WORKDIR /loj/judge

COPY go.mod go.sum /loj/judge/

RUN go mod download -x

COPY ./ /loj/judge

RUN go generate ./cmd/executorserver/version \
    && CGO_ENABLE=0 go build -v -tags nomsgpack -o loj_judger_server ./cmd/executorserver

FROM gcc

WORKDIR /opt

COPY --from=build /loj/judge/loj_judger_server  /loj/judge/mount.yaml /opt/

EXPOSE 6060/tcp

ENTRYPOINT ["./loj_judger_server"]
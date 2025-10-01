FROM registry.cn-shenzhen.aliyuncs.com/juxuny-public/builder:golang-1.23.10 AS builder
COPY go.mod /src/
COPY go.sum /src/
RUN cd /src && go env -w GOPROXY=https://goproxy.cn && go mod download
COPY . /src/
RUN go env
RUN cd /src && go mod tidy && CGO_ENABLED=0 go build -o app

# final stage
FROM registry.cn-shenzhen.aliyuncs.com/juxuny-public/builder:alpine-3.8-tz
WORKDIR /app
RUN cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && echo 'Asia/Shanghai' > /etc/timezone
COPY --from=builder /src/app /app/entrypoint
ENTRYPOINT /app/entrypoint -l :20000 -r 127.0.0.1:10000 --to-proxy-protocol --from-de-xun

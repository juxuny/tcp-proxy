FROM golang:1.16.4 AS builder
#ARG GOPRIVATE
#ARG GIT_USER
#ARG GIT_TOKEN
#RUN git config --global url."https://$GIT_USER:$GIT_TOKEN@$GOPRIVATE".insteadOf "https://$GOPRIVATE"
COPY go.mod /src/
COPY go.sum /src/
RUN cd /src && go env -w GOPROXY=https://goproxy.cn && go mod download
COPY . /src/
RUN go env
RUN cd /src && \
    CGO_ENABLED=0 go build -o app

# final stage
#FROM golang:1.16.4
FROM ineva/alpine:3.9
WORKDIR /app
RUN cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && echo 'Asia/Shanghai' > /etc/timezone
COPY --from=builder /src/app /app/entrypoint
ENTRYPOINT /app/entrypoint -l :20000 -r 127.0.0.1:10000 --to-proxy-protocol --from-de-xun

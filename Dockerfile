
FROM golang:1.21.3 As Yumo

#创建工作区目录
RUN mkdir /app

# 将本地代码复制到镜像
COPY ./ /app

#指定工作区目录
WORKDIR /app

RUN go mod tidy

#打包二进制文件
RUN go build main.go

#拉取精简镜像
FORM scratch

COPY --from=Yumo /app/main /main

CMD ["./main.go"]
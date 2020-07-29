FROM golang:alpine
WORKDIR /app
COPY . /app
ENV GOPROXY https://goproxy.cn
ENTRYPOINT ["ash", "-c", "go clean && go build -o /env/tree-hole tree-hole && cd /env && ./tree-hole"]

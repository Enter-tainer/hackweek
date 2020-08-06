# tree-hole NO.1037 后端

文档见 `doc` 目录。

## 部署

更改 `env/production.json`，然后运行

```
docker-compose start
```

服务会在 3000 端口上启动。

## 设计

使用 golang 编写，使用了 mongodb 和 redis。使用 jwt 来进行用户鉴权。

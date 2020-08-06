# tree-hole NO.1037 后端

文档见 `doc` 目录。

## 部署

更改 `env/production.json`，然后运行

```
docker-compose start
```

服务会在 3000 端口上启动。

## 设计

使用 golang 编写，使用 http 框架 echo，数据库采用了 mongodb 和 redis，并使用 jwt 来进行用户鉴权。

### 匿名性

- 密码加盐，数据库中存储 argon2 哈希。
- 所有 user_id 仅在一个帖子内有效，不会泄露用户的真实 id。

# 数据库结构文档

## Mongo

### 用户集合 user

- _id `ObjectId`：用户 ID
- password `String`：bcrypt hash 后的密码
- email `String`：邮箱
- verified `Boolean`：是否已经验证邮箱

## Redis

- `code:<邮箱验证码>`：验证码对应的用户 ID

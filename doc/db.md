# 数据库结构文档

## Mongo

### 用户集合 user

- _id `ObjectId`：用户 ID
- password `String`：argon2 hash 后的密码
- email `String`：邮箱
- verified `Boolean`：是否已经验证邮箱

### 帖子集合 post

- _id `ObjectID`：帖子 ID
- title `String` 标题
- createdAt `Date` 发布时间
- updatedAt `Date` 修改时间
- user `ObjectID` 发布人ID
- content `String`: 帖子内容
- reply `Array`：所有回复

reply 格式与 posts 类似，但是没有 title 和 reply 两个字段

## Redis

- `code:<邮箱验证码>`：验证码对应的用户 ID

# 最简单的 CRUD Demo 

## 使用 Vue 3.3 + Vite + js + Element-Plus 实现




项目技术

前端 :

Vue3.2 + Vite + ElementPlus + js

后端 :

Go + Gin + GORM + Mysql

所需环境

Node.js
Go
Mysql

已经实现的功能

前端 CRUD
后端 CRUD



Vue3.2 + Go + Axios 数据联调

yarn add axios



创建 MySQL 数据库
```sql
create schema crud collate utf8mb4_0900_ai_ci;

```

```shell
# 启动前端
yarn dev

# 启动后端服务
go run main.go
```


浏览器访问：

http://localhost:5173/
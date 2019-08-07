# [WIP] db-querier
一个简单的指定SQL查询web服务，支持多数据源，多数据库

> - 后端使用Go语言，[gin](https://github.com/gin-gonic/gin)配合[sqlx](https://github.com/jmoiron/sqlx)开发
> - 前端使用[vue.js](https://vuejs.org/)

#
1. 在后端的配置文件中动态配置好要执行的query以及制定参数（配置支持热加载）
2. 前端页面会动态渲染出支持的查询列表和参数列表，并有基本的类型检查

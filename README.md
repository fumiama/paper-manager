# paper-manager
面向教师的试卷网络管理系统，基于Golang+vben。

## 项目结构
项目由前端后端两大部分组成，没有中间件。
### frontend
基于`vben`编写前端页面并编译为静态文件以供`main.go`嵌入。
### backend
基于标准`http`库，不用重框架，负责提供各种API。
### main.go
统合前后端为单一可执行文件，静态访问交给前端，API部分则转交后端处理。

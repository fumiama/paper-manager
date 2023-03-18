# paper-manager
面向教师的试卷网络管理系统，基于Golang+vben。

# 项目结构
项目由前端后端两大部分组成，没有中间件。
## frontend
基于`vben`编写前端页面并编译为静态文件以供`main.go`嵌入。
### 登录页
> /login

#### 登录
输入账号密码登录。登录成功后，前端将查询并缓存用户名、权限、头像、简介等信息备用，同时导航到指定的家页面。
- **课程组长**：导航到分析页`/dashboard/analysis`
- **其他人**：导航到工作台`/dashboard/workbench`
#### 忘记密码
点击`忘记密码`后填写用户名与手机号码，再点击`重置`，即可将重置消息报告给课程组长，由课程组长电话联系确认无误后，在系统中进行密码重置。
#### 注册
点击`注册`后填写用户名、手机号与密码，再点击`注册`，即可将注册消息报告给课程组长，由课程组长电话联系确认无误后，在系统中创建相应账号。
### 仪表板/分析页
> /dashboard/analysis

### 仪表板/工作台
> /dashboard/workbench

### 个人设置
> /settings

个人设置位于右上角状态栏头像的下拉列表中，点击后即可对用户自己的信息进行设置。
#### 基本设置
对用户的昵称、个人简介与头像进行自定义设置。
#### 安全设置
对用户的密码、联系方式进行修改。
- **密码**：导航至`/settings/password`进行设置，成功后将自动退出当前登录，同时在系统消息中通知课程组长。
- **联系方式**：导航至`/settings/contact`进行设置，成功后会在系统消息中通知课程组长。

### 试卷库
> /filelist


## backend
基于标准`http`库，不用重框架，负责提供各种API。
## main.go
统合前后端为单一可执行文件，静态访问交给前端，API部分则转交后端处理。
### 参数
- **-l**: 设置监听地址与端口
### 监听处理点位
- **/api/\*\*\***: 所有信令 API
- **/file/\*\*\*xxx**: 所有上传的文件（动态文件）
- **/upload**: 上传文件接口，上传的文件会统一存至`./data/file/`以供访问

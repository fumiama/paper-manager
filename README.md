# paper-manager
面向教师的试卷网络管理系统，基于Golang+vben。项目由前端后端两大部分组成，没有中间件。

# main.go
统合前后端为单一可执行文件，静态访问交给前端，API部分则转交后端处理。
## 参数
- **-l**: 设置监听地址与端口
## 监听处理点位
- **/api/\*\*\***: 所有信令 API
- **/file/\*\*\***: 所有上传的文件（动态文件）
- **/upload**: 上传文件接口，上传的文件会统一存至`./data/file/`以供访问


# 前端
> 位于`frontend`文件夹

基于`vben`编写前端页面并编译为静态文件以供`main.go`嵌入。

## 登录页 /login

### 登录
输入账号密码登录。登录成功后，前端将查询并缓存用户名、权限、头像、简介等信息备用，同时导航到指定的家页面。
- **课程组长**：导航到分析页`/dashboard/analysis`
- **其他人**：导航到工作台`/dashboard/workbench`
> 登录时将依次访问`/api/getLoginSalt` `/api/login` `/api/getUserInfo`

![login screen](https://user-images.githubusercontent.com/41315874/226117983-c1e69916-def0-4746-939a-5041412b755f.png)

### 忘记密码
点击`忘记密码`后填写用户名与手机号码，再点击`重置`，即可将重置消息报告给课程组长，由课程组长电话联系确认无误后，在系统中批准进行密码重置。
> 上报消息时访问`/api/resetPassword`

<table>
	<tr>
		<td align="center"><img src="https://user-images.githubusercontent.com/41315874/226824943-4b8f8a24-4393-4c90-bd27-f7ca38e8f01c.png"></td>
		<td align="center"><img src="https://user-images.githubusercontent.com/41315874/226824951-ae61494c-d6ce-4358-8fa3-f6b24eee3afc.png"></td>
	</tr>
    <tr>
		<td align="center">填写信息</td>
		<td align="center">上报消息</td>
	</tr>
</table>

### 注册
点击`注册`后填写用户名、手机号与密码，再点击`注册`，即可将注册消息报告给课程组长，由课程组长电话联系确认无误后，在系统中批准以所填信息创建账号。
> 上报消息时访问`/api/register`

<table>
	<tr>
		<td align="center"><img src="https://user-images.githubusercontent.com/41315874/226804117-40cfe055-d4ca-440f-9129-1599a126f520.png"></td>
		<td align="center"><img src="https://user-images.githubusercontent.com/41315874/226804339-936e38df-e8ea-4d85-a29f-75f191538fee.png"></td>
	</tr>
    <tr>
		<td align="center">填写信息</td>
		<td align="center">上报消息</td>
	</tr>
</table>

## 个人设置 /settings

<img align="right" src="https://user-images.githubusercontent.com/41315874/226120888-ce79c227-7f0f-4681-ab74-99ce4433d768.png" alt="settings" />

个人设置位于右上角状态栏头像的下拉列表中，点击后即可对用户自己的信息进行设置。
### 基本设置
对用户的昵称、个人简介与头像进行自定义设置。
> 上传头像时访问`/upload`

> 保存设置时访问`/api/setUserInfo`

![base setting](https://user-images.githubusercontent.com/41315874/226120541-e2cc77e7-6601-49bb-8f8e-fe6d348f210b.png)

### 安全设置
对用户的密码、联系方式进行修改。
- **密码**：导航至`/settings/password`进行设置，成功后将自动退出当前登录，同时在系统消息中通知课程组长。
- **联系方式**：导航至`/settings/contact`进行设置，成功后会在系统消息中通知课程组长。
> 设置密码时访问`/api/setPassword`

> 设置联系方式时访问`/api/setContact`

![secure setting](https://user-images.githubusercontent.com/41315874/226120560-83a8fa95-f2db-4202-8aee-a99c32b55b43.png)

## 注销

<img align="right" src="https://user-images.githubusercontent.com/41315874/226120865-9f8d57bf-3884-420e-9ff6-008f50fb52d6.png" alt="logout" />

该功能位于右上角状态栏头像的下拉列表中，点击`退出系统`后即可注销登录。
> 注销时访问`/api/logout`

## 仪表板/分析页 /dashboard/analysis
此页仅课程组长可见，向课程组长展示近一年的访问量信息。
> 载入时访问`/api/getAnnualVisits`

![analysis](https://user-images.githubusercontent.com/41315874/226802920-f6b43a7b-6191-4dcb-9f48-364c161c1cfd.png)

## 仪表板/工作台 /dashboard/workbench
显示用户收到的消息，用户可以选择`删除`已读消息。对于课程组长，还提供对`注册`和`重置密码`消息的快捷处理按钮`接受`，可以快速执行相应操作。下面是一个演示。
> 载入时访问`/api/getUsersCount` `/api/getMessageList`

> 处理消息时访问`/api/acceptMessage` `/api/delMessage`
### 1. 课程组长工作台
> 登入系统，修改了自己的信息，又收到了两个用户的注册请求

![workbench demo 1](https://user-images.githubusercontent.com/41315874/226808521-92b2a857-4b56-4185-b654-759e2d43d415.png)

### 2. 课程组长工作台
> 接受2人注册请求

![workbench demo 2](https://user-images.githubusercontent.com/41315874/226808582-ebea1c32-a69f-4054-8b46-8f3f5799eb9f.png)

### 3. A老师(filemgr)工作台
> 课程组长接受注册申请后，更改了个人信息

![workbench demo 3](https://user-images.githubusercontent.com/41315874/226806112-344bafc1-70f1-4408-bcd0-2cd6896c7269.png)

### 4. B老师(user)工作台
> 课程组长接受注册申请后，更改了个人信息和联系方式

![workbench demo 4](https://user-images.githubusercontent.com/41315874/226806143-ce5a833a-f30d-4f51-945e-bf4b83459758.png)

### 5. 课程组长工作台
> 最终的消息状态

![workbench demo 5](https://user-images.githubusercontent.com/41315874/226811656-a6e4238f-2dbe-4a49-a50f-3d2075838cb8.png)

## 仪表板/用户管理 /dashboard/account
此页仅课程组长可见，提供用户管理功能，可以禁用用户、编辑用户权限和信息。下面是一个管理用户的例子。
> 载入时访问`/api/getUsersList`

> 禁用账户时访问`/api/disableUser`

> 设置用户信息时访问`/api/setUserInfo`

> 设置用户权限时访问`/api/setRole`

### 1. 上一板块演示后的用户管理界面状态
![account demo 1](https://user-images.githubusercontent.com/41315874/226825676-e943d841-eb5f-4292-b4d5-a7f91cf48962.png)
### 2. 课程组长编辑filemgr用户的信息
![account demo 2](https://user-images.githubusercontent.com/41315874/226825963-0f4c70aa-d1ec-4215-a7ed-d81c540a66d4.png)
### 3. A老师(filemgr)工作台上的显示
![account demo 3](https://user-images.githubusercontent.com/41315874/226826317-daa537a7-abaf-44b0-9377-d057573223cd.png)
### 4. 课程组长禁用filemgr用户
![account demo 4](https://user-images.githubusercontent.com/41315874/226826425-66a42bcf-5a5a-455c-8aad-1117b7c3ccae.png)
### 5. filemgr用户无法登陆
![account demo 5](https://user-images.githubusercontent.com/41315874/226826516-cfdc18c7-f59f-4759-a008-7baca34f1a30.png)
### 6. filemgr用户通过重置密码向课程组长发出恢复账号申请
![account demo 6](https://user-images.githubusercontent.com/41315874/226826700-99e8a639-58db-499f-9629-756fc5e5e49f.png)
### 7. 课程组长接受申请
![account demo 7](https://user-images.githubusercontent.com/41315874/226826829-6dec3a7a-8386-4d33-aa6e-cea26cc846e9.png)
### 8. A老师(filemgr)成功登录进入工作台
![account demo 8](https://user-images.githubusercontent.com/41315874/226826932-2c365095-a35e-477f-ace7-a98452c24db4.png)

## 试卷库 /filelist


# 后端
> 位于`backend`文件夹

基于标准`http`库，不用重框架，负责提供各种API。

## 当前进展及测试情况
```
已完成:
A.log模块基础已经好了 -- ok
B.测试-f传参其他路径的文件是否有效，拿viper取值测试即可，两个路径的value不同 -- ok
C.router模块配置 --ok
D.DB模块配置  -- 等待测试同步建库 -- 当前:需要建表、测试同步、测试api结果 -- ok
E.将所有的实例化的函数中加个判断，如果是nil再实例化，否则直接返回--解决程序重复实例化，浪费资源问题 -- ok
F.panic+recover -- ok
G.尝试添加runtime.Caller()方法到error模块中。方便快速拿到报错信息的行号 -- ok
H.登录认证，失败多次则无法登录，一段时间后可以登录。 -- ok
I.解决gorm的自定义logger(gorm.SetLogger()) -- 有时间要了解一下
J.jwt用户token授权及认证 -- ok
K.k8s操作 -- ok
L.chartmuseum操作 -- ok
M.tiller操作 -- ok
O.用户权限表先弄好:
  1)用户权限表创建: user role为主  -- ok
  2)自动创建admin用户  -- ok
P.增加用户访问提示功能: -- ok
server/server.go InitServer()函数 --> middleware/access_log.go
s.addMiddleWare(
		middleware.AccessLog(skipLogPath...),
提示效果: INFO[2021-01-08 16:35:11] ip: 10.151.150.13   latency: 141.509µs  code: 500   method: POST     path: /v1/auth/login


未完成:
A.加swagger,所有的validate要加example -- 
B.把main.go中的一堆flag都加上，方便命令传参 -- 
C.在弄用户操作: 
  1)基于admin用户进行普通用户的创建(用户名不可以有特殊符号)，并发送邮件 -- ok
  2)普通用户 login、logout
  3)admin用户 login、logout
  4)修改账户信息(密码、邮箱)
D.增加用明文生成salt和加密password的功能，方便后期更新admin用户密码。
E.drone-agent和drone-server的使用方法(README.md)、脚本(传参、启动服务)、文件(二进制文件、脚本文件等)都要备份到github工程中 *****

当前进展:
        0. 所有的delete操作，先进行查询，报错是new出来的。方便postman获取错误id时有正常的报错显示。
        1. 完善v1/user.go的 所有部分增加 序列化 MarshalResponse

1.验证 tokenExpire没有配置的后果 -- ok (时间设置为 5m 10m)
2.验证 token是否为设置的5分钟，不要因为计算错误，按照秒数来算的分钟就完蛋了
  需要登录后才有效果，暂时只有login生成token，还要有其他api，才能验证过期的效果
3.加一下deleted_at都models/common中 -- ok
4.测试deleted_at字段是否好用
5.尝试不用Localtime，created_at updated_at用默认的time.Time
6.解决role外键问题 -- ok,但不建议使用，所以关闭了
7.确认报错，这是有外键束缚的时候的报错，后面不做束缚了，但得确认报错
(E:/mygopath/src/go.rock.com/rock-platform/rock/server/database/api/user.go:37)
[2021-01-06 23:10:15]  [1.97ms]  INSERT INTO `user` (`name`,`password`,`email`,`salt`,`token`,`login_retry_count`,`login_block_until`,`role_id`,`created_at`,`updated_at`,`deleted_at`,`version`) VALUES ('rock4','f0cfbcfff9b52ae091a1e72b81bafb44','123456@qq.com','uvFqtTrATP8H3Hnx','',0,NULL,0,'2021-01-06 23:10:15','2021-01-06 23:10:15',NULL,0)
[0 rows affected or returned ]
ERRO[2021-01-06 23:10:15] Mysql Error: with num 1452 and message is: Cannot add or update a child row: a foreign key constraint fails (`rock`.`user`, CONSTRAINT `user_role_id_role_id_foreign` FOREIGN KEY (`role_id`) REFERENCES `role` (`id`))
8.cluster admin.conf 用postman上传的时候是一行内容，保存到数据库应该是多行才对，需要检测下。 -- ok 正常的
9.models里面的version字段 作用是什么，为啥一直是0呢？
```

## 未知知识
```
drone plugins
```

## console与rock代码行数
```
console:
# for i in `find . -type f -name '*.go'  | grep -v './vendor/' | egrep -v 'clients/license'` ;do cat $i | wc -l  | tr "\n" "+"; done
# 额外再加 ./server/clients/license/cactl.go 585行
142+145+9552+199+33+96+26+78+34+93+101+43+79+170+115+49+53+129+26+41+39+451+155+183+346+206+209+14+22+273+239+564+631+404+324+243+244+223+201+98+319+198+392+233+158+72+175+182+202+163+90+232+213+193+410+334+107+80+42+18+17+42+36+46+63+13+29+61+22+24+69+21+104+79+61+64+140+18+25+170+78+27+163+14+28+48+43+171+87
共计: 21,849行

Rock platform:
46+139+1238+25+13+241+16+18+3+1+242+254+108+137+285+101+19+55+29+24+29+112+77+60+51+138+18+85+65+9+45+70+93+104
共计: 3,950行(2020-12-28 11:30:05)  // 准确
46+139+1968+43+39+25+13+237+241+16+18+3+187+82+245+321+141+108+111+143+318+104+22+19+55+29+17+24+29+112+77+60+60+138+18+18+113+66+9+45+72+93+104
共计: 5,728行(2021-02-07 11:56:16)
46+139+2577+40+37+25+13+222+241+16+18+3+230+90+245+321+137+108+145+143+318+106+22+19+55+29+17+24+29+112+77+60+60+138+18+18+144+66+9+45+72+91+104
共计: 6,429行(2021-02-08 21:50:03)
共计: 7,366行(2021-02-19 14:15:40)
共计: 10,765行(2021-02-23 22:06:04) // 准确
共计: 12,516行(2021-03-04 20:39:08) // 准确
共计: 13,889行(2021-03-10 20:36:03) // 准确
共计: 16,101行(2021-03-18 21:03:26) // 准确
共计: 17,846行(2021-04-06 15:35:55) // 准确(不含clients/license)
```

## 知识点疑问
+ 1 `drone-agent`
  - 1.1 `drone-agent`源码位置在哪里?
  - 1.2 `drone-agent`的作用主要是什么? 是`drone-server`发来的任务，按照pipeline去执行任务吗?
  - 1.3 如果`.drone.yaml`任务都是`drone-agent`执行，那么是如何启动容器的呢?
  
+ 2 `drone-server`
  - 2.1 是否是用于`drone-go`模块与`drone-server`二进制交互，通过`access-token`(`drone-token`)获取该用户的gitlab相关信息、接收任务的增删改查请求。
  - 2.2 `drone-server`拿到增删改查任务后，发给`drone-agent`进行处理吗?

+ 3.Rock理解发版逻辑
  - 前端进行单个服务发版，将相关参数传给CreateApp Api，
  - CreateApp Api 通过drone-go模块将相关参数发给 drone-server，
  - drone-server将任务下发给drone-agent，
  - drone-agent 拉取该应用的源码，根据 .drone.yaml(pipeline)定义进行任务执行，

## 全量验证-最终必须做的事情
+ 1.使用自己的`drone-server`源码编译成二进制，进行验证使用。
+ 2.使用自己的`drone-agent`源码编译成二进制，进行验证使用。
+ 3.更新自己的`infra-drone-plugins`，找到镜像源，制作成镜像，并保存，最后使用这个镜像进行部署验证。
+ 4.使用自己的`infra-frontend-service`源码编译成镜像及chart，进行验证使用。

## 收尾工作
```
1. 状态码更新好
2. 报错还是默认走就行，要么logger报错，要么panic，不要再弄 runtime.Caller 了。
3. 优化user.go等里面的返回值，改为v1.xxx。
4. 增加及配置 debug日志，方便后期排查。
5. 转为vendor保存
```

## 遗留难解问题
```
1.login登录时，用不存在的账号登录，postman接收的错误为空{}，但日志报错是 record not found。需要解决postman为空问题。
2.deleted_at 无法加到models中，因为总是报错，不知道tag该填什么值 -- 已修复，用time.Time类型即可
3.一些应该warning的报错，不应该eror报出来:
   ERRO[2021-01-07 10:18:49] Rock Error: func_name: [doPrintf], file_name:[print.go], line:[1030] :user with name(rock3) is alerady exist
4. 用顺义定义的RoleID *RoleIdReq 方式报错，必须int64才不报错
E:\mygopath\src\go.rock.com\rock-platform\rock\server\controller\v1\user.go 
48 if err := ctx.ShouldBind(&userReq); err != nil {
//RoleId   *RoleIdReq `json:"role_id" binding:"required"`  // 用顺义的这种定义，ctx.ShouldBind报错
抱错内容: json: cannot unmarshal number into Go value of type v1.RoleIdReq 
为48行ctx.ShouldBind(&userReq)的err报错
```
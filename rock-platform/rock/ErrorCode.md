### 错误码对照表
```
const(
	StatusContinue           = 100
	StatusSwitchingProtocols = 101

	StatusOK                   = 200
	StatusCreated              = 201
	StatusAccepted             = 202
	StatusNonAuthoritativeInfo = 203
	StatusNoContent            = 204
	StatusResetContent         = 205
	StatusPartialContent       = 206

	StatusMultipleChoices   = 300
	StatusMovedPermanently  = 301
	StatusFound             = 302
	StatusSeeOther          = 303
	StatusNotModified       = 304
	StatusUseProxy          = 305
	StatusTemporaryRedirect = 307

	StatusBadRequest                   = 400
	StatusUnauthorized                 = 401
	StatusPaymentRequired              = 402
	StatusForbidden                    = 403
	StatusNotFound                     = 404
	StatusMethodNotAllowed             = 405
	StatusNotAcceptable                = 406
	StatusProxyAuthRequired            = 407
	StatusRequestTimeout               = 408
	StatusConflict                     = 409
	StatusGone                         = 410
	StatusLengthRequired               = 411
	StatusPreconditionFailed           = 412
	StatusRequestEntityTooLarge        = 413
	StatusRequestURITooLong            = 414
	StatusUnsupportedMediaType         = 415
	StatusRequestedRangeNotSatisfiable = 416
	StatusExpectationFailed            = 417
	StatusTeapot                       = 418
    StatusUnprocessableEntity          = 422

	StatusInternalServerError     = 500
	StatusNotImplemented          = 501
	StatusBadGateway              = 502
	StatusServiceUnavailable      = 503
	StatusGatewayTimeout          = 504
	StatusHTTPVersionNotSupported = 505

	// New HTTP status codes from RFC 6585. Not exported yet in Go 1.1.
	// See discussion at https://codereview.appspot.com/7678043/
	statusPreconditionRequired          = 428
	statusTooManyRequests               = 429
	statusRequestHeaderFieldsTooLarge   = 431
	statusNetworkAuthenticationRequired = 511
)

常用状态码:
201 StatusCreated 创建成功
400 StatusBadRequest 错误请求
401 StatusUnauthorized 权限不足
404 StatusNotFound 请求不存在
500 StatusInternalServerError 服务器错误
```

### 400错误类型
- HttpCode: 400
- ErrorCode: 40000001
- 中文名称: 注册用户名已存在

---
- HttpCode: 400
- ErrorCode: 40000002
- 中文名称: 密码长度错误

---
- HttpCode: 400
- ErrorCode: 40000003
- 中文名称: 密码错误

---

- HttpCode: 400
- ErrorCode: 40000004
- 中文名称: 用户名不存在

---

- HttpCode: 400
- ErrorCode: 40000006
- 中文名称: 序列化失败

---

- HttpCode: 400
- ErrorCode: 40000007
- 中文名称: 密码强度不够

--- 

- HttpCode: 400
- ErrorCode: 40000008
- 中文名称: 新密码不能和老密码相同

---

- HttpCode: 400
- ErrorCode: 40000009
- 中文名称: 该邮箱已注册

---

- HttpCode: 400
- ErrorCode: 40000010
- 中文名称: 两次密码输入不同

---

- HttpCode: 400
- ErrorCode: 40000011
- 中文名称: 该邮箱不存在

---

- HttpCode: 400
- ErrorCode: 40000012
- 中文名称: secret超时

---

- HttpCode: 400
- ErrorCode: 40000013
- 中文名称: secret错误

---
- HttpCode: 400
- ErrorCode: 40000015
- 中文名称: project项目名称已经存在

---

- HttpCode: 400
- ErrorCode: 40000017
- 中文名称: app应用名称已经存在

---

- HttpCode: 400
- ErrorCode: 40000018
- 中文名称: cluster集群名称已经存在

---

- HttpCode: 400
- ErrorCode: 40000019
- 中文名称: 传入的k8s config错误

---

- HttpCode: 400
- ErrorCode: 40000020
- 中文名称: 该cluster_id已存在对应namespace

---

- HttpCode: 400
- ErrorCode: 40000021
- 中文名称: admin用户的role_id不允许修改

---

- HttpCode: 400
- ErrorCode: 40000022
- 中文名称: 该project_id env_id name已存在

---

- HttpCode: 400
- ErrorCode: 40000023
- 中文名称: Annotation console.cluster.id不能转为int类型

---

- HttpCode: 400
- ErrorCode: 40000024
- 中文名称: app config不是正确的yaml格式

---

- HttpCode: 400
- ErrorCode: 40000025
- 中文名称: k8s config不是正确的yaml格式

---

- HttpCode: 401
- ErrorCode: 40100001
- 中文名称: 权限不足

---

- HttpCode: 403
- ErrorCode: 40300001
- 中文名称: 关联的Gitlab账号，对该工程没有master权限

---

- HttpCode: 404
- ErrorCode: 40400001
- 中文名称: cookie中没有发现config

---

- HttpCode: 404
- ErrorCode: 40400002
- 中文名称: admin用户禁止删除

---

- HttpCode: 404
- ErrorCode: 40400003
- 中文名称: app应用id不存在

---

- HttpCode: 404
- ErrorCode: 40400004
- 中文名称: project项目id不存在

---

- HttpCode: 404
- ErrorCode: 40400005
- 中文名称: 用户id不存在

---

- HttpCode: 404
- ErrorCode: 40400006
- 中文名称: role_id不存在

---

- HttpCode: 404
- ErrorCode: 40400007
- 中文名称: cluster_id不存在

---

- HttpCode: 404
- ErrorCode: 40400008
- 中文名称: env_id不存在

---

- HttpCode: 404
- ErrorCode: 40400009
- 中文名称: project_env_id不存在

---

- HttpCode: 404
- ErrorCode: 40400010
- 中文名称: 此app应用gitlab_project_id不存在

---

- HttpCode: 404
- ErrorCode: 40400011
- 中文名称: deployment_id不存在

---

- HttpCode: 404
- ErrorCode: 40400012
- 中文名称: instance_id不存在

---

- HttpCode: 404
- ErrorCode: 40400013
- 中文名称: 当前程序不支持这个版本的K8S

---

- HttpCode: 404
- ErrorCode: 40400014
- 中文名称: 不存在该app的appConf配置文件

---

- HttpCode: 404
- ErrorCode: 40400015
- 中文名称: k8s集群节点不存在

---

- HttpCode: 412
- ErrorCode: 41200001
- 中文名称: 缺少gitlab access token

---

- HttpCode: 412
- ErrorCode: 41200002
- 中文名称: 该project_env_id存在app_conf config,请删该config后再删除该project_env

### 500错误类型
---
- HttpCode: 500
- ErrorCode: 50000001
- 中文名称:  服务器发生未知错，请联系管理员处理

---
- HttpCode: 500
- ErrorCode: 50000002
- 中文名称: 未知数据库错误

---
- HttpCode: 500
- ErrorCode: 50000003
- 中文名称: 未知邮件系统错误

---
- HttpCode: 500
- ErrorCode: 50000004
- 中文名称: 未知Drone错误

---
- HttpCode: 500
- ErrorCode: 50000005
- 中文名称: tiller端口ping不通

---
- HttpCode: 500
- ErrorCode: 50000006
- 中文名称: helm安装chart失败


```
- HttpCode: 400
- ErrorCode: 40000001
- 中文名称: 您发送的请求参数或格式不符合要求，请您确认相关请求接口的规范

---
- HttpCode: 400
- ErrorCode: 40000002
- 中文名称: 您填写的集群配置信息无法通过验证，请您确认相关集群配置

---
- HttpCode: 400
- ErrorCode: 40000003
- 中文名称: 您创建的资源已经存在

---
- HttpCode: 400
- ErrorCode: 40000004
- 中文名称: 您与认证服务交互的证书不正确

---
- HttpCode: 400
- ErrorCode: 40000005
- 中文名称: 您应该调用的周期任务接口

---
- HttpCode: 400
- ErrorCode: 40000006
- 中文名称: 您应该调用的任务接口

---
- HttpCode: 400
- ErrorCode: 40000007
- 中文名称: 您提供的V2C文件中可激活的次数已经用尽

---
- HttpCode: 400
- ErrorCode: 40000008
- 中文名称: 您对权限和角色执行的操作动作不被支持

---
- HttpCode: 400
- ErrorCode: 40000009
- 中文名称: 您提供的原密码不正确

---
- HttpCode: 400
- ErrorCode: 40000010
- 中文名称: 您无法修改超级管理员的角色属性

---
- HttpCode: 400
- ErrorCode: 40000011
- 中文名称: 您提供的链接已经失效，请重新获取新的链接后再试

---
- HttpCode: 400
- ErrorCode: 40000012
- 中文名称: 您提供的链接不正确，请重新获取新的链接后再试

---
- HttpCode: 400
- ErrorCode: 40000013
- 中文名称: 您对用户和角色执行的操作动作不被支持

---
- HttpCode: 400
- ErrorCode: 40000014
- 中文名称: 您解除的某些权限不属于当前的角色

---
- HttpCode: 400
- ErrorCode: 40000015
- 中文名称: 您使用的集群配置被损坏，无法提取集群地址信息

---
- HttpCode: 400
- ErrorCode: 40000016
- 中文名称: 无法解析cookie中的配置信息

---
- HttpCode: 400
- ErrorCode: 40000017
- 中文名称: 无法解析cookie中的配置信息

---
- HttpCode: 400
- ErrorCode: 40000018
- 中文名称: 您提供的YAML格式的文件有语法错误，请修改后重试

---
- HttpCode: 400
- ErrorCode: 40000019
- 中文名称: 您无法删除超级管理员用户

---
- HttpCode: 401
- ErrorCode: 40100001
- 中文名称: 您没有登录,请登录

---
- HttpCode: 401
- ErrorCode: 40100002
- 中文名称: 您的登录信息无效，请重新登录

---
- HttpCode: 401
- ErrorCode: 40100003
- 中文名称: 您的登录信息已变更，请重新登录

---
- HttpCode: 401
- ErrorCode: 40100004
- 中文名称: 您登录的用户名或密码错误，请重新登录

---
- HttpCode: 401
- ErrorCode: 40100005
- 中文名称: 您尝试登陆失败次数已达上限，请稍后重试

---
- HttpCode: 403
- ErrorCode: 40300001
- 中文名称: 您没有权限，所以操作被禁止

---
- HttpCode: 403
- ErrorCode: 40300002
- 中文名称: 您关联的Gitlab账号，对该工程没有master权限，所以操作被禁止

---
- HttpCode: 404
- ErrorCode: 40400001
- 中文名称: 您查看的资源或者相关联的资源不存在，请重新确认

---
- HttpCode: 404
- ErrorCode: 40400002
- 中文名称: 您的cookie中无法找到配置信息,请确认开启cookie

---
- HttpCode: 404
- ErrorCode: 40400003
- 中文名称: 您的应用没有关联代码仓库的具体工程，无法执行该操作

---
- HttpCode: 412
- ErrorCode: 41200001
- 中文名称: 您正在删除的资源有其他资源依赖，无法删除，请您先删除所有的依赖项

---
- HttpCode: 412
- ErrorCode: 41200002
- 中文名称: 您需要提前设置gitlab的access token

---
- HttpCode: 412
- ErrorCode: 41200003
- 中文名称: 您删除的应用有关联的服务实例存在，请先删除所有关联的服务实例

---
- HttpCode: 412
- ErrorCode: 41200004
- 中文名称: 您删除的集群有关联的资源空间存在，请先删除所有关联的资源空间

---
- HttpCode: 412
- ErrorCode: 41200005
- 中文名称: 您删除的资源空间有关联的服务实例存在，请先删除所有关联的服务实例

---
- HttpCode: 412
- ErrorCode: 41200006
- 中文名称: 您删除的权限有关联的角色存在，请先删除所有关联的角色

---
- HttpCode: 412
- ErrorCode: 41200007
- 中文名称: 您删除的项目有关联的应用存在，请先删除所有关联的应用

---
- HttpCode: 412
- ErrorCode: 41200008
- 中文名称: 您删除的项目空间有关联的应用配置存在，请先删除关联的应用配置

---
- HttpCode: 412
- ErrorCode: 41200009
- 中文名称: 您删除的角色有关联的用户存在，请先解除所有关联关系

---
- HttpCode: 412
- ErrorCode: 41200010
- 中文名称: 您删除的角色有关联的权限存在，请先解除所有关联关系

---
- HttpCode: 412
- ErrorCode: 41200011
- 中文名称: 您删除的用户有关联的角色存在，请先解除所有关联关系
---


### 500错误类型
- HttpCode: 500
- ErrorCode: 50000001
- 中文名称: 服务器发生未知错，请联系管理员处理

---
- HttpCode: 500
- ErrorCode: 50000002
- 中文名称: 在线激活操作失败

---
- HttpCode: 500
- ErrorCode: 50000003
- 中文名称: 离线激活操作失败

---
- HttpCode: 500
- ErrorCode: 50000004
- 中文名称: 离线激活操作执行后，证书依然处于未激活状态

---
- HttpCode: 500
- ErrorCode: 50000005
- 中文名称: 未知数据库错误

---
- HttpCode: 500
- ErrorCode: 50000006
- 中文名称: 未知kubernetes错误

---
---
- HttpCode: 500
- ErrorCode: 50000007
- 中文名称: 未知Tiller状态,连接错误

---
- HttpCode: 500
- ErrorCode: 50000008
- 中文名称: 在线激活操作执行后，证书依然处于未激活状态

---     
```
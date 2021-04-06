### 1 drone官网文档
```
http://docs.drone.io/step-conditions/ --> https://docs.drone.io/pipeline/docker/syntax/conditions/
http://docs.drone.io/pipelines/
```
### 2 drone源码位置
```
drone-agent  源码位置: 
drone-server 源码位置: https://gitlab.sz.sensetime.com/galaxias/drone.git 分支为master
drone-go 源码位置: https://github.com/rockwang465/drone/drone-go (go代码引用此模块和drone-server交互)
```

### 3 `.drone.yaml`理解
#### `.drone.yaml`示例
```yaml
workspace:
  base: /app
  path: ./

environment:
  REGISTRY: 10.151.3.75
  CHARTMUSEUM: http://10.151.3.75:8080
  API_SERVICE: http://10.151.3.85:32000
  REPO_NAME: idea-aurora/aurora-push-service  # 10.151.3.75 harbor镜像仓库中的前缀
  APP_VERSION: 1.0.0
  CHART_PLUGIN: 10.151.3.75/cicd/infra-cd  # chart打包的镜像(自己制作)
  DOCKER_PLUGIN: 10.151.3.75/cicd/docker-plugin  # drone官网提供的镜像(plugins/docker-> http://plugins.drone.io/drone-plugins/drone-docker/)，未做任何修改
  DEPLOY_PLUGIN: 10.151.3.75/galaxias/infra-drone-plugins  # 部署到环境的镜像(自己制作)

pipeline:

  build-and-push-image:  # 编译代码生成镜像并推到harbor仓库
    image: ${DOCKER_PLUGIN}
    registry: ${REGISTRY}
    secrets: [ docker_username, docker_password ]
    repo: ${REGISTRY}/${REPO_NAME}
    insecure: true
    tags:
      - "${APP_VERSION}-${DRONE_BRANCH}-${DRONE_COMMIT_SHA:0:6}"  # DRONE_BRANCH DRONE_COMMIT_SHA 通过drone进行传参的

  package_and_upload_chart:  # 打包生成chart版本包并推到chartmuseum仓库
    image: ${CHART_PLUGIN}
    chart_image_name: ${REGISTRY}/${REPO_NAME}
    chart_image_tag: "${APP_VERSION}-${DRONE_BRANCH}-${DRONE_COMMIT_SHA:0:6}"
    repo_addr: ${CHARTMUSEUM}

  deploy_to_env:
    image: 10.151.3.75/galaxias/infra-drone-plugins  # 使用该镜像启动一个容器
    secrets: [ galaxias_api_token ]  # 通过galaxias进行传参
    environment:
      CHART_REPO_NAME: aurora-push-service
      CHART_REPO_VERSION: "${APP_VERSION}-${DRONE_BRANCH}-${DRONE_COMMIT_SHA:0:6}"
      DEBUG: 'True'
    commands:
      - /plugins/deploy.py  # 调用运维平台 /v1/deployments 接口进行指定chart版本服务更新(chart版本号来自)
```
#### 3.1 `.drone.yaml`详解
##### 3.1.1 `environment`
+ 预先定义好变量，方便后面pipeline部分引用变量

##### 3.1.2 `build-and-push-image`
```yaml
    image: ${DOCKER_PLUGIN}  # docker run镜像
    registry: ${REGISTRY}
    secrets: [ docker_username, docker_password ]
    repo: ${REGISTRY}/${REPO_NAME}
    insecure: true
    tags:
      - "${APP_VERSION}-${DRONE_BRANCH}-${DRONE_COMMIT_SHA:0:6}"  # DRONE_BRANCH DRONE_COMMIT_SHA 通过drone进行传参的
```
+ `image: 10.151.3.75/cicd/docker-plugin`
  - A.`image`: 表示使用此镜像启动容器
  - B.`10.151.3.75/cicd/docker-plugin`: 此为drone官网提供的镜像(`plugins/docker-> http://plugins.drone.io/drone-plugins/drone-docker/`)，未做任何修改
  - C.drone-docker源码地址: `https://github.com/drone-plugins/drone-docker`
  - D.drone-docker启动命令: 

+ `registry: ${REGISTRY}`
  


#### 3.2 疑问解答
##### 3.2.1 `deploy.py`脚本源码位于何处?
```
答: https://gitlab.sz.sensetime.com/galaxias/infra-drone-plugins/blob/master/plugins/deploy.py
    此脚本已经制作成镜像,用于 deploy_to_env 中使用此镜像,进行服务的部署
```
##### 3.2.2 `deploy.py`脚本做了什么?
```
答: 通过galaxias和.drone.yaml拿到一堆环境变量,如:cookies(galaxias_api_token(jwt token))/docker_username/docker_password/app_id/project_env_id/env_id/chart_name/chart_version等信息
    见3.4 启动容器，执行deploy.py, 传入上面的环境变量,向 http://10.151.3.xx:8888/v1/deployments(当前运维平台后端) api 发起请求, 部署指定版本的chart tgz包
```
##### 3.2.3 `galaxias_api_token`如何生成与作用?
```
答: 生成方法:
    galaxias_api_token 是 galaxias平台的admin账号的jwt加密后的token,用于 deploy.py 中的请求cookies参数进行认证.
    而这个token保存在数据库的user表中,所以可以到这里来拿.
    10.151.3.85(galaixias)的admin token: token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJ1c2VybmFtZSI6ImFkbWluIiwiZHJvbmVfdG9rZW4iOiIiLCJwYXNzd29yZCI6IjU0NThiZmU4OWU5N2RjZjc3OTk5MDEwNGI0ZjU2YTU4Iiwicm9sZSI6ImFkbWluIiwiaXNzIjoiY29uc29sZSJ9.yTo1-t_5vXsd0Ywbv447DOrfHAP3JXeuSWsq2wQccTw
    10.151.3.86(rock)的admin      token: token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJ1c2VybmFtZSI6ImFkbWluIiwiZHJvbmVfdG9rZW4iOiIiLCJwYXNzd29yZCI6IjMyMDdlYWQ0ZTA5MmRlNzdlMDIyMzk0YjMyMDRkNzU1Iiwicm9sZSI6ImFkbWluIiwiZXhwIjoxNjE0MDc3NjAzLCJpYXQiOjE2MTQwNzE1NDMsImlzcyI6IlJvY2sgV2FuZyIsInN1YiI6IkxvZ2luIHRva2VuIn0.nAMR3xjGZ-4etgyVT2qfiUx2oEZhKM_iRs8lui1vTJ4
    作用: 从 deploy.py脚本中理解, galaxias_api_token是admin的jwt加密后的token. 用于 deploy.py 中的请求cookies参数进行认证.
```
##### 3.2.4 `.drone.yaml`中是如何执行 `/plugins/deploy.py` 脚本的?
```
答: .drone.yaml 到 deploy_to_env 这一步的时候, 会起image为 infra-drone-plugins 的镜像, 然后进入容器执行 /plugins/deploy.py 这条command来执行此脚本.
```

### 4 drone-agent 的作用,galaxias是如何连接到drone-agent的?  -- 后期明白逻辑后再更新一下这里的逻辑
```
答: 作用: 作为drone的agent端,接收drone-server发来的任务, 然后基于pipeline(.drone.yaml)执行任务.
    示例: 需要部署一个chart服务,前端点击构建某个服务的任务(参数为:app_id、branch/tag、环境ip、名称空间等),
         后端go代码将通过drone-go模块交互将构建任务发给drone-server(二进制),
         drone-server将任务下发给drone-agent, drone-agent拿到.drone.yaml的pipeline内容进行工作:代码克隆、docker build(代码编译及构建image)、helm pack和使用deploy.py脚本进行部署到指定环境中。
```

### 5 `start_drone_dev_agent.sh`脚本中 `DRONE_SECRET=ZHJvbmUtbXlzcWwK` 这个`secret`的作用是什么，如何生成的?
```
```
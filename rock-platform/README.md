# Rock Platform 岩石平台
## 1.功能介绍
+ 前提条件:需要将操作的环境信息录入到数据库
+ 对集群环境进行信息展示(集群的节点信息、集群的k8s信息)
+ 对集群环境k8s的底层服务进行重启
+ 对集群环境chart服务进行版本展示、更新、删除、override配置文件更新
+ 对chartmuseum仓库的版本展示

## 2.底层组件
### 1)mysql(gorm)
+ 用于各节点信息的数据存储(用户token、admin.conf)

## 3.使用前准备
### 3.1 配置`config.yaml`到`/etc/rock`目录下
### 3.2 创建数据库`rock`
```text
# docker run mysql
10.151.3.85
缺少本地bin-log数据缓存，实际项目请注意
docker run --restart=always --name mysql-test -p 3333:3306 -e MYSQL\_ROOT\_PASSWORD=rock1314 -d mysql:5.7.23
mysql -uroot -P3333 -h0.0.0.0 -prock1314
mysql> create database rock character set UTF8mb4 collate utf8mb4_bin; 
mysql> ALTER USER 'root'@'%' IDENTIFIED WITH mysql_native_password BY 'rock1314';

# add iptables rules
# iptables -I INPUT -p tcp -m tcp --dport 3333 -j ACCEPT
# iptables -I INPUT -m state --state ESTABLISHED,RELATED -j ACCEPT
```
### 3.3 创建数据库`drone`
+ docker run mysql for drone
```text
缺少本地bin-log数据缓存，实际项目请注意
# docker run --restart=always --name mysql-drone -p 3306:3306 -e MYSQL\_ROOT\_PASSWORD=123456 -d mysql:5.7.23
```

+ create database
```text
# docker exec -it mysql-drone bash
# mysql -uroot -p123456
mysql> CREATE USER 'drone'@'localhost'  IDENTIFIED BY '123456';  # 创建用户
mysql> CREATE USER 'drone'@'%'  IDENTIFIED BY '123456';  # 远程登录
mysql> GRANT ALL PRIVILEGES ON *.* to 'drone'@'%' identified by '123456';  # 授权用户所有权限
mysql> CREATE DATABASE drone55 character set UTF8mb4 collate utf8mb4_bin;  # 建库
mysql> ALTER USER 'root'@'%' IDENTIFIED WITH mysql_native_password BY '123456';
mysql> ALTER USER 'drone'@'%' IDENTIFIED WITH mysql_native_password BY 'drone_123456';
mysql> FLUSH PRIVILEGES;
```

+ add iptables rules
```text
# iptables -I INPUT -p tcp -m tcp --dport 3306 -j ACCEPT
# iptables -I INPUT -m state --state ESTABLISHED,RELATED -j ACCEPT
```

### 3.4 解决不能使用`GROUP BY`问题
```text
mysql> SELECT @@sql_mode;
+-------------------------------------------------------------------------------------------------------------------------------------------+
| @@sql_mode                                                                                                                                |
+-------------------------------------------------------------------------------------------------------------------------------------------+
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION |
+-------------------------------------------------------------------------------------------------------------------------------------------+
1 row in set (0.00 sec)

对新建数据库有效:
mysql> SET @@global.sql_mode ='STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION';

对于已存在的数据库需要执行:
mysql> SET sql_mode ='STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION';
```

## 4.使用介绍
### 1)启动命令
+ `server` 启动服务
+ `-v` 查看版本号
+ `--config /etc/xxx/config.yaml` 指定配置文件
+ `--log-dir /xxx/log/` 指定日志路径
```
# ./rock  server -h
start rock server

Usage:
  rock server [flags]

Flags:
  -c, --config string    The config file's path (default "/etc/rock/config.yml")
  -h, --help             help for server
      --log-dir string   The log file's directory (default "/var/log/rock")
```

### 2)swagger
+ swagger页面
```
http://10.151.3.87:8888/swagger/index.html
```
+ 初始化
```
# cd /e/mygopath/src/github.com/swaggo/swag/cmd/swag
# cd /e/mygopath/src/go.rock.com/rock-platform/rock/cmd
# /e/mygopath/src/github.com/swaggo/swag/cmd/swag/swag.exe init -g ./cmd/main.go - o ./docs
2020/12/20 14:51:32 Generate swagger docs....
2020/12/20 14:51:32 Generate general API Info, search dir:./
2020/12/20 14:51:33 warning: failed to get package name in dir: ./, error: execute go list command, exit status 1, stdout:, stderr:can't load package: package .: no Go files in E:\mygopath\src\go.rock.com\rock-platform\rock
2020/12/20 14:51:33 create docs.go at docs\docs.go
2020/12/20 14:51:33 create swagger.json at docs\swagger.json
2020/12/20 14:51:33 create swagger.yaml at docs\swagger.yaml
```

### 3)跨平台编译
```
# cd /e/mygopath/src/rock-platform/rock/cmd
# CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build .
```

## 使用`infra-frontend`chart配合二进制`rock-platform`
+ 创建service关联endpoint
```yaml
# vim infra-console-service-svc.yaml
apiVersion: v1
kind: Service
metadata:
  name: infra-console-service
  namespace: dev
spec:
  ports:
  - protocol: TCP
    port: 8888
    targetPort: 8888
  clusterIP: None
---
apiVersion: v1
kind: Endpoints
metadata:
  name: infra-console-service
  namespace: dev
subsets:
  - addresses:
      - ip: 10.151.3.86
    ports:
      - port: 8888
```
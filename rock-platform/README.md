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
+ 配置`config.yaml`到`/etc/rock`目录下
+ 创建数据库`rock`
```
docker run mysql
10.151.3.85
docker run --restart=always --name mysql-test -p 3333:3306 -e MYSQL\_ROOT\_PASSWORD=rock1314 -d mysql:5.7.23
mysql -uroot -P3333 -h0.0.0.0 -prock1314
mysql> create database rock character set UTF8mb4 collate utf8mb4_bin; 
mysql> ALTER USER 'root'@'%' IDENTIFIED WITH mysql_native_password BY 'rock1314';
# iptables -I INPUT -p tcp -m tcp --dport 3333 -j ACCEPT
# iptables -I INPUT -m state --state ESTABLISHED,RELATED -j ACCEPT
```

## 4.使用介绍
### 1)启动命令
+ `server` 启动服务
+ `-v` 查看版本号
+ `--config /etc/xxx/config.yaml` 指定配置文件
+ `--log-dir /xxx/log/` 指定日志路径

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

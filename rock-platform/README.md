# Rock Platform 岩石平台
## 功能介绍
+ 前提条件:需要将操作的环境信息录入到数据库
+ 对集群环境进行信息展示(集群的节点信息、集群的k8s信息)
+ 对集群环境k8s的底层服务进行重启
+ 对集群环境chart服务进行版本展示、更新、删除、override配置文件更新
+ 对chartmuseum仓库的版本展示

## 底层组件
### mysql(gorm)
+ 用于各节点信息的数据存储(用户token、admin.conf)

## 使用介绍
+ `-f` 指定配置文件
+ `-v` 查看版本号
+ `server` 启动服务

## 跨平台编译
```
# cd /e/mygopath/src/rock-platform/rock/cmd
# CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build .
```

## `docker run mysql`
```
docker run --name mysql-test -p 3333:3306 -e MYSQL\_ROOT\_PASSWORD=rock1314 -d mysql
mysql -uroot -P3333 -h0.0.0.0 -prock1314
mysql> create database demo2;
```

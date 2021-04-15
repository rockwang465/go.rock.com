#!/bin/bash

# solve push code error: unable to access 'http://github.com/xxxx.git' OpenSSL SSL_read: SSL_ERROR_SYSCALL, errno 10054
# 方法一:
# git config --global http.postBuffer 524288000  # 修改上传大小上限
# git config --global http.sslVerify "false"  # ssl验证关闭
# 方法二(测试好用):
# 开始菜单 -> 更改代理设置 -> 自动检测设置 -> 关掉(之前开启的)

echo -e "go build to ./rock :"
#CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./rock .
#CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./rock -mod vendor .
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./rock -mod vendor go.rock.com/rock-platform/rock/cmd .

echo -e "scp rock to 10.151.3.86:/rock"
scp rock root@10.151.3.86:/rock

echo -e `date`

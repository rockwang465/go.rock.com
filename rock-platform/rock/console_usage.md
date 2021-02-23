## infra-console-service使用介绍
### 1.服务基础信息
+ 服务启动在`10.151.3.85`机器的`k8s`的`dev`名称空间上
+ 配置文件放在`configmap`中
+ 数据库是`10.151.3.85`的`3306`端口`console55`数据库
+ drone脚本位置:`/root/drone-test/infra/dev/start_drone_dev_agent.sh`和`/root/drone-test/infra/dev/start_drone_dev_server.sh`

### 2.数据库
+ 数据库创建
```sql
CREATE DATABASE `console` /*!40100 COLLATE 'utf8_unicode_ci' */;
```
```sql
CREATE USER 'console'@'%' IDENTIFIED BY 'UVlY88m9suHLsthK';
GRANT ALL PRIVILEGES ON console.* to console@'%' IDENTIFIED by 'UVlY88m9suHLsthK';
FLUSH PRIVILEGES;
```

### 3.配置文件
+ config配置文件
```yaml
# vim /etc/console/config.yml
insecure-bind-address: 0.0.0.0
insecure-port: 8080
cookie-expire: 6h
secret-expire: 12h
login-retry-count: 3
login-block-duration: 5m
license-ca-port: 32043
drone:
  addr: http://10.151.3.99:8888
repo:
  addr: http://10.151.3.99:38080
log:
  level: debug
  dir: "."
db:
  driver: mysql
  user: console
  password: rock1314
  address: 10.151.3.85
  port: 3333
  name: console
tiller:
  port: 31134
email:
  user: wangyecheng465@163.com
  password: OXQBIWAMOOAYOHNV
  smtp:
    addr: smtp.163.com
    port: 465
frontend:
  domain: http://localhost:8989
```

### 4.参数
+ 传参用法
```text
# ./console run -h
Show how to run infra console server

Usage:
  console run [flags]

Flags:
  -c, --config string                  The config file's path (default "/etc/console/config.yml")
      --cookie-expire string           The expire time duration of cookie, valid unit is: h(hour), m(minute). e.x. 1.5h, 30m, 72h (default "6h")
      --db-address string              The DB's address (default "127.0.0.1")
      --db-driver string               The DB driver for this server (default "mysql")
      --db-name string                 The DB's database name (default "console")
      --db-password string             The DB's password
      --db-port string                 The DB's port (default "3306")
      --db-user string                 The DB's user (default "console")
      --drone-addr string              Drone server's address (default "http://127.0.0.1:8000")
      --email-name string              The system admin email user who would send system email to all users
      --email-password string          The system admin email user's password
      --email-smtp-addr string         The email SMTP server addr (default "smtp.partner.outlook.cn")
      --email-smtp-port int            The email SMTP server port (default 587)
      --frontend-domain string         The frontend domain, e.x. http://galaxias.sensetime.com:8080
  -h, --help                           help for run
      --insecure-bind-address string   The server insecure http address where listen on
      --insecure-port string           The server insecure http port where listen on
      --license-ca-port int            The license CA server port (default 8443)
      --log-dir string                 The log file's directory (default "/var/log/console")
      --log-level string               The log level,options: debug,info,warn,error (default "debug")
      --login-block-duration string    The block duration when user login failed count come up to password retry count, valid unit is: h(hour), m(minute). e.x. 1.5h, 30m, 72h (default "5m")
      --login-retry-count int          The failed login count to trigger block user login (default 3)
      --repo-addr string               Charts museum repo's domain or ip address (default "http://127.0.0.1:8080")
      --secret-expire string           The secret valid time duration, valid unit is: h(hour), m(minute). e.x. 1.5h, 30m, 72h (default "24h")
  -b, --secure-bind-address string     The server secure https address where listen on (default "127.0.0.1")
  -p, --secure-port string             The server secure https port where listen on (default "443")
      --tiller-port int                The tiller's port (default 44134)
      --tls-cert-file string           File containing the default x509 Certificate for HTTPS
      --tls-private-key-file string    File containing the default x509 private key matching --tls-cert-file
```

### 5.脚本
+ 执行脚本
```text
# cat execute_console.sh  # 支持文件、传参方式
./console run
./console run -c="/etc/console/config.yml" --db-address="10.151.3.85" --db-password="rock1314" --email-name="wangyecheng465@163.com" --email-password="OXQBIWAMOOAYOHNV"  --email-smtp-addr="smtp.163.com" --email-smtp-port="465"
```
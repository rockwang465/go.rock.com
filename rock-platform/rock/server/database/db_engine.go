package database

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"go.rock.com/rock-platform/rock/server/conf"
	"go.rock.com/rock-platform/rock/server/database/models"
	"go.rock.com/rock-platform/rock/server/log"
	"net/url"
	"strings"
)

type DBEngine struct {
	*gorm.DB
}

var SingleDBEngine *DBEngine

type DBSource struct {
	Driver   string
	User     string
	Password string
	Host     string
	Port     int
	Name     string
	Charset  string
	Loc      string
}

func ConnectDB() *gorm.DB {
	config := conf.GetConfig()
	var dbSource = &DBSource{
		Driver:   config.Viper.GetString("db.driver"),
		User:     config.Viper.GetString("db.user"),
		Password: config.Viper.GetString("db.password"),
		Host:     config.Viper.GetString("db.host"),
		Port:     config.Viper.GetInt("db.port"),
		Name:     config.Viper.GetString("db.name"),
		Charset:  config.Viper.GetString("db.charset"),
		Loc:      url.QueryEscape(config.Viper.GetString("db.loc")),
	}
	// "root:rock1314@tcp(10.151.3.86:3333)/rock?charset=utf8mb4&parseTime=True&loc=Asia%2FShanghai"
	args := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=%s",
		dbSource.User,
		dbSource.Password,
		dbSource.Host,
		dbSource.Port,
		dbSource.Name,
		dbSource.Charset,
		dbSource.Loc)
	db, err := gorm.Open(dbSource.Driver, args)
	//fmt.Printf("%v\n", args)
	if err != nil {
		panic(err)
	}
	db.SingularTable(true) // 禁止表名复数
	return db
}

func GetDBEngine() *DBEngine {
	if SingleDBEngine == nil {
		dbEngine := ConnectDB()
		SingleDBEngine = &DBEngine{dbEngine}
	}
	return SingleDBEngine
}

// close DB
func Close() {
	if err := SingleDBEngine.DB.Close(); err != nil {
		panic(err)
	}
}

func (e *DBEngine) InitDB() {
	config := conf.GetConfig()
	logger := log.GetLogger()

	// set gorm log level
	logLevel := config.Viper.GetString("log.level")
	switch strings.ToLower(logLevel) {
	case "debug": // debug模式需要展示sql语句，所以这里单独初始化了日志等级
		e.LogMode(true) // 会在控制台打印出生成的SQL语句
	case "info", "warn", "error":
		e.LogMode(false)
	default: // 没有匹配上则默认debug
		e.LogMode(true)
		logger.Warningf("[Rock Platform] Got unknown log level %s, and set DB log level to default: debug", logLevel)
	}
	logger.Infof("[Rock Platform] Set DB's log level: %s", logLevel)

	// sync tables
	logger.Infof("[Rock Platform] Start to sync tables ...")
	var user = &models.User{} // 或者直接 user := new(models.User),赋值直接用&user
	var role = &models.Role{}
	var project = &models.Project{}
	var app = &models.App{}
	var cluster = &models.Cluster{}
	var env = &models.Env{}
	e.AutoMigrate(role, user) // create role user table
	e.AutoMigrate(project, app, cluster, env)
	//e.Model(user).AddForeignKey("role_id", "role(id)", "RESTRICT", "RESTRICT") // add ForeignKey,正常业务不应该开启外键束缚,影响数据库性能

	logger.Infof("[Rock Platform] Tables sync finished")
}

package log

import (
	"github.com/gin-gonic/gin"
	"github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
	"go.rock.com/rock-platform/rock/server/conf"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Logger struct {
	*logrus.Logger
}

var SingleLogger *Logger

func GetLogger() *Logger {
	if SingleLogger == nil {
		SingleLogger = &Logger{logrus.New()}
		formatter := &logrus.TextFormatter{
			ForceColors:     true,
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
		}
		SingleLogger.Logger.SetFormatter(formatter)
	}

	return SingleLogger
}

func (l *Logger) InitLogger() {
	config := conf.GetConfig()

	// Set log directory location
	logDir := config.Viper.GetString("log.dir")
	l.Infoln("[Rock Platform] Set log dir to: ", logDir)

	// Set rotation log options  日志切割
	logFile := filepath.Join(logDir, "rock.%Y%m%d.log")
	logf, err := rotatelogs.New(logFile,
		rotatelogs.WithLinkName("rock.log"),       // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(7*time.Hour*24),     // 文件最大保存时间
		rotatelogs.WithRotationTime(time.Hour*24), // 日志切割时间间隔
	)
	if err != nil {
		l.Fatalf("[Rock Platform] failed to create rotation logs: %s", err)
	}
	gin.DefaultWriter = io.MultiWriter(logf, os.Stderr) // gin.DefaultWriter变量能控制日志的保存方式及保存位置
	log.SetOutput(gin.DefaultWriter)                    // 决定了log应该输出到什么地方，默认是标准输出（同下解释）
	l.SetOutput(gin.DefaultWriter)                      // 设置output,默认为stderr,可以为任何io.Writer，比如文件*os.File

	//Set log level
	logLevel := config.Viper.GetString("log.level")
	switch strings.ToLower(logLevel) {
	case "debug":
		l.SetLevel(logrus.DebugLevel) // Debug级别应该会打到屏幕上的
	case "info":
		l.SetLevel(logrus.InfoLevel)
		gin.SetMode(gin.ReleaseMode) // 全局设置环境，gin.DebugMode为开发环境，gin.ReleaseMode线上环境为
	case "warn":
		l.SetLevel(logrus.WarnLevel)
		gin.SetMode(gin.ReleaseMode)
	case "error":
		l.SetLevel(logrus.ErrorLevel)
		gin.SetMode(gin.ReleaseMode)
	default:
		l.SetLevel(logrus.DebugLevel)
		l.Warningf("[Rock Platform] Got unknown log level %s, and set log level to default: debug", logLevel)
	}
	l.Infoln("[Rock Platform] Set log level to:", logLevel)
}

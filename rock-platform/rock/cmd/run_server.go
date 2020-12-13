package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"net/http"
	"path"
	"rock-platform/rock/server"
	"rock-platform/rock/server/conf"
	"rock-platform/rock/server/utils"
	"strings"
)

func runServer(cmd *cobra.Command, args []string) {
	config := conf.GetConfig()
	// https://www.liwenzhou.com/posts/Go/viper_tutorial/ 了解一下3个参数用法
	// 自动读取环境变量参数。BindEnv和AutomaticEnv函数会使用SetEnvPrefix中定义的前缀。
	config.Viper.AutomaticEnv()
	// 通过设置环境变量前缀SetEnvPrefix,在从环境变量读取时会添加设置的前缀。来确保 Env 变量是唯一的(小写会自动转为大写)
	config.Viper.SetEnvPrefix("ROCK")
	// 批量替换: 将 -. 都转为 _
	replacer := strings.NewReplacer("-", "_", ".", "_")
	// 设置环境变量的分隔符
	config.Viper.SetEnvKeyReplacer(replacer)

	// Load config file by viper
	configFile := config.Viper.GetString("config")            // /etc/rock/config.yaml
	fileDir := path.Dir(configFile)                           //  文件目录 /etc/rock
	fileFormat := strings.TrimLeft(path.Ext(configFile), ".") // 获取文件类型 yaml
	fileName := utils.GetConfigName(path.Base(configFile))    // 获取文件名的头(去掉.yaml): config
	config.Viper.SetConfigType(fileFormat)
	config.Viper.SetConfigName(fileName)
	config.Viper.AddConfigPath(fileDir) // 第一个搜索路径
	config.Viper.AddConfigPath(".")     // 可以多次调用添加路径，比如添加当前目录
	err := config.Viper.ReadInConfig()  // Find and read the config file

	sv := server.GetServer() // 实例化: logrus、gin路由、xorm数据库

	if err != nil {
		sv.Logger.Fatalf("[Rock Platform] Fatal error on reading config faile: %s \n", err)
	} else {
		sv.Logger.Infoln("[Rock Platform] Got config file path set to:", configFile)
	}

	sv.InitServer() // 初始化日志配置、中间件、路由、数据库、validator

	// http server configuration
	httpAddr := config.Viper.GetString("server.addr")
	httpPort := config.Viper.GetInt64("server.port")

	httpServer := &http.Server{}
	// 为空serveHttp为false，不为空serveHttp为true
	serverHttp := httpAddr != "" && httpPort != 0
	if serverHttp {
		listen := fmt.Sprintf("%s:%d", httpAddr, httpPort)
		httpServer.Addr = listen
		httpServer.Handler = sv.RouterEngine
		sv.Logger.Infoln("[Rock Platform] Set http listen addr to: ", listen)

		// server.ListenAndServer()
		go func() {
			if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				sv.Logger.Fatalf("[Rock Platform] listen http address error, and process exit: %s", err)
			}
		}()
	} else { // 不使用http，则要用https了
		sv.Logger.Infoln("[Rock Platform] Server skip http listen set up")
	}

	sv.Logger.Infoln("Welcome use Rock Platform")

	// Just test viper
	serverPort := config.Viper.GetInt64("server.port")
	if serverPort == 0 {
		serverPort = 8000 // set default port
	}
	sv.Logger.Infoln("Rock platform server port is :", serverPort)
}

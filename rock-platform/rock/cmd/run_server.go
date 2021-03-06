package main

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	_ "go.rock.com/rock-platform/rock/docs" // docs is generated by Swag CLI, you have to import it.
	"go.rock.com/rock-platform/rock/server"
	"go.rock.com/rock-platform/rock/server/conf"
	"go.rock.com/rock-platform/rock/server/utils"
	"net/http"
	"os"
	"os/signal"
	"path"
	"strings"
	"time"
)

// @title Swagger API Doc
// @version 1.0
// @description This is a api doc for this server.
// @termsOfService https://swaggo.github.io/swaggo.io/declarative_comments_format/
func runServer(cmd *cobra.Command, args []string) {

	config := conf.GetConfig()
	// https://www.liwenzhou.com/posts/Go/viper_tutorial/ 了解一下3个参数用法
	// 自动读取环境变量参数。BindEnv和AutomaticEnv函数会使用SetEnvPrefix中定义的前缀。
	config.Viper.AutomaticEnv()
	// 通过设置环境变量前缀SetEnvPrefix,在从环境变量读取时会添加设置的前缀。来确保 Env 变量是唯一的(小写会自动转为大写)
	config.Viper.SetEnvPrefix("CONSOLE")
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

	sv := server.GetServer() // 实例化: logrus、gin路由、gorm数据库

	if err != nil {
		sv.Logger.Fatalf("[Rock Platform] Fatal error on reading config file: %s \n", err)
	} else {
		sv.Logger.Infoln("[Rock Platform] Got config file path set to:", configFile)
	}

	sv.InitServer() // 初始化日志配置、中间件、路由、数据库、validator

	// http server configuration
	httpAddr := config.Viper.GetString("server.addr")
	httpPort := config.Viper.GetInt64("server.port")

	if httpPort == 0 {
		httpPort = 8000 // define a default port
	}

	httpServer := &http.Server{}
	// 为空serveHttp为false，不为空serveHttp为true
	serveHttp := httpAddr != "" && httpPort != 0
	if serveHttp {
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

	// https server configuration
	httpsAddr := config.Viper.GetString("server.https.addr")
	httpsPort := config.Viper.GetInt64("server.https.port")
	tlsCertFile := config.Viper.GetString("server.https.tls-cert-file")
	tlsPrivateKeyFile := config.Viper.GetString("server.https.tls-private-key-file")

	httpsServer := &http.Server{}
	serveHttps := tlsCertFile != "" && tlsPrivateKeyFile != ""
	if serveHttps {
		listen := fmt.Sprintf("%s:%d", httpsAddr, httpsPort)
		httpsServer.Addr = listen
		httpsServer.Handler = sv.RouterEngine
		sv.Logger.Infoln("[Rock Platform] Set https listen addr to:", listen)

		// server.ListenAndServer()
		go func() {
			if err := httpsServer.ListenAndServeTLS(tlsCertFile, tlsPrivateKeyFile); err != nil && err != http.ErrServerClosed {
				sv.Logger.Fatalf("[Rock Platform] listen https address error, and process exit, %s", err)
			}
		}()
	} else {
		sv.Logger.Infoln("[Rock Platform] Server skip https listen set up")
	}

	if !serveHttp && !serveHttps {
		sv.Logger.Fatalln("[Rock Platform] You must set up at least one of http or https")
		return
	}

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit // 正常此处会阻塞,当收到中断程序的信号时，就会拿到内容，执行下面的代码

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if serveHttp { // 停止服务提示
		sv.Logger.Println("[Rock Platform] Shutting down HTTP Server ...")
		if err := httpServer.Shutdown(ctx); err != nil {
			sv.Logger.Fatal("[Rock Platform] Http server Shutdown:", err)
			return
		}
	}

	if serveHttps {
		sv.Logger.Println("[Rock Platform] Shutting down HTTPS Server ...")
		if err := httpsServer.Shutdown(ctx); err != nil {
			sv.Logger.Fatal("[Rock Platform] Https server Shutdown:", err)
			return
		}
	}

	sv.Logger.Println("[Rock Platform] Server exited successfully")
	return
}

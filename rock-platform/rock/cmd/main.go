package main

import (
	"github.com/spf13/cobra"
	"go.rock.com/rock-platform/rock/server/conf"
	"go.rock.com/rock-platform/rock/server/log"
)

func main() {
	var cmdCfg = conf.GetCmdCfg() // 配置文件信息
	var config = conf.GetConfig() // viper实例化
	var logger = log.GetLogger()  // logrus实例化

	var rootCmd = &cobra.Command{
		Use:     "rock [command]",
		Short:   "This is a operation platform",
		Example: "rock command",
		Version: "v1.0.0",
	}

	var runServerCmd = &cobra.Command{
		Use:   "server",
		Short: "start rock server",
		Run:   runServer,
		Args:  cobra.ArbitraryArgs, // 命令接受任何参数，且不返回错误
	}

	// cmdCfg.ConfigPath 用于把value: "/etc/console/config.yml"中的值赋值过去
	// 传参用法: -c 或 --config = "/etc/rock/config.yml"
	runServerCmd.Flags().StringVarP(&cmdCfg.ConfigPath, "config", "c", "/etc/rock/config.yml",
		"The config file's path")
	// 将subCmd.Flags()命令行参数绑定到viper中，解决万一有用户命令行传参，则可以替换默认读取的值。
	_ = config.Viper.BindPFlag("config", runServerCmd.Flags().Lookup("config"))

	// 传参用法: --log-dir="/var/log/rock"
	runServerCmd.Flags().StringVar(&cmdCfg.LogDir, "log-dir", "/var/log/rock",
		"The log file's directory")
	_ = config.Viper.BindPFlag("log.dir", runServerCmd.Flags().Lookup("log-dir"))

	rootCmd.AddCommand(runServerCmd)
	err := rootCmd.Execute()
	if err != nil {
		logger.Logger.Errorln("Error : root cmd execute err : ", err)
		panic(err)
	}
}

/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"tarantula-v2/app"
	"tarantula-v2/config"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "tarantula-v2",
	Short: "通过Amazon 商品详情页面获取商品信息截图",
	Long:  `执行命令将启动Rabbit 客户端通过监听指定消息队列来获取品类信息的查询请求. 例如：`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		// 初始化配置
		initGlobalVariables()
		// 启动消费者
		go app.Consuming(viper.GetString("mq.url"), viper.GetString("mq.exchange"),
			viper.GetString("mq.exchangeType"), viper.GetString("mq.queue"))

		fmt.Println("已城通启动消费者，监听中...")
		// 永不退出
		select {}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", ".tarantula.yaml", "config file (default is $HOME/.tarantula.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".tarantula-v2" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".tarantula-v2")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

// 初始化全局变量
func initGlobalVariables() {
	// 是否启用全局浏览器上下文
	config.GlobalContext = viper.GetBool("global-context")
	if config.GlobalContext {
		fmt.Println("启用全局浏览器上下文，将不会在每次请求时重新创建浏览器上下文，请求结束后也不会关闭浏览器上下文。")
		err := config.InitBrowserContext("")
		if err != nil {
			fmt.Println("初始化浏览器上下文失败，请检查浏览器远程端口是否打开...", err)
			return
		}
	}
}

package conf

import (
	"ecloud_computer_auto_boot/pkg/util"
	"errors"
	"github.com/spf13/viper"
	"os"
)

type server struct {
	Debug    bool   `yaml:"debug" mapstructure:"debug"`
	LogLevel string `yaml:"log-level" mapstructure:"log-level"`
}

type secret struct {
	Type      string `yaml:"type" mapstructure:"type"`
	Username  string `yaml:"username" mapstructure:"username"`
	Password  string `yaml:"password" mapstructure:"password"`
	AccessKey string `yaml:"access-key" mapstructure:"access-key"`
	SecretKey string `yaml:"secret-key" mapstructure:"access-key"`
	PoolId    string `yaml:"pool-id" mapstructure:"pool-id"`
}

type cron struct {
	Duration int      `yaml:"duration" mapstructure:"duration"`
	Machines []string `yaml:"machines" mapstructure:"machines"`
}

func Init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath(".")

	defaultConfig := map[string]any{
		"server": Server,
		"secret": Secret,
		"cron":   Cron,
	}

	for key, val := range defaultConfig {
		viper.SetDefault(key, val)
	}

	if err := viper.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			err := viper.SafeWriteConfig()
			if err != nil {
				util.Log().Error("[配置初始化] 默认配置文件写出失败: %s", err)
				os.Exit(1)
			}

			util.Log().Info("[配置初始化] 已生成默认配置文件，请配置完成后再次运行。")
			os.Exit(0)
		}

		util.Log().Error("[配置初始化] 配置文件读取失败: %s", err)
		os.Exit(1)
	}

	for key, val := range defaultConfig {
		err := viper.UnmarshalKey(key, val)
		if err != nil {
			util.Log().Error("[配置初始化] 配置文件解析失败, key: %s, error: %s", key, err)
		}
	}

	if Server.Debug {
		Server.LogLevel = "debug"
	}

	// 重设log等级
	if Server.LogLevel != "" {
		util.GlobalLogger = nil
		util.BuildLogger(Server.LogLevel)
	}
}

package configs

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var defaultConfigurations = map[string]string{
	awsRegion: "eu-north-1",

	savingPath: "./file_sync/",

	syncPeriodSeconds: "15",
	filterTags:        "",
	filterPrefixName:  "",
}

func Init() {
	viper.SetEnvPrefix("AWS_SECRETS_STORAGE_SYNC")
	loadConfig()
}

func loadConfig() {
	loadDefaultConfigs()
	loadFileConfigs()
	loadEnvConfigs()
}

func loadDefaultConfigs() {
	for configKey, configValue := range defaultConfigurations {
		viper.SetDefault(configKey, configValue)
	}
}

func loadFileConfigs() {
	viper.SetConfigName("config")
	viper.AddConfigPath("./configs")
	err := viper.ReadInConfig()
	if err != nil {
		logrus.Warn("can't load config from file. Use Variable environments and Default")
	}
}

func loadEnvConfigs() {
	viper.AutomaticEnv()
}

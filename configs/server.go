package configs

import (
	"fmt"

	"github.com/spf13/viper"
)

const (
	serverIP   = "SERVER_IP"
	serverPort = "SERVER_PORT"
)

func GetServerIP() string {
	return viper.GetString(serverIP)
}

func GetServerPort() int {
	return viper.GetInt(serverPort)
}

func GetServerAddress() string {
	return fmt.Sprintf("%s:%d", GetServerIP(), GetServerPort())
}

package configs

import "github.com/spf13/viper"

const (
	awsRegion = "AWS_REGION"
)

func GetAWSRegion() string {
	return viper.GetString(awsRegion)
}

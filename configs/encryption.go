package configs

import "github.com/spf13/viper"

const (
	awsEncryptionMethod         = "ENCRYPTION_METHOD"
	EncryptionAES256Base64Key   = "AES_GCM_256_BASE_64_KEY"
	EncryptionAES256Base64Nonce = "AES_GCM_256_BASE_64_NONCE"
)

type EncryptionMethod string

const (
	EncryptionMethodNone   EncryptionMethod = "NONE"
	EncryptionMethodAES256 EncryptionMethod = "AES_GCM_256"
)

type EncryptionConfig struct {
	Method EncryptionMethod
	Key    string
	Nonce  string
}

func GetEncryptionConfig() EncryptionConfig {
	return EncryptionConfig{
		Method: GetEncryptionMethod(),
		Key:    GetEncryptionAES256Base64Key(),
		Nonce:  GetEncryptionAES256Base64Nonce(),
	}
}

func GetEncryptionMethod() EncryptionMethod {
	return EncryptionMethod(viper.GetString(awsEncryptionMethod))
}

func GetEncryptionAES256Base64Key() string {
	return viper.GetString(EncryptionAES256Base64Key)
}

func GetEncryptionAES256Base64Nonce() string {
	return viper.GetString(EncryptionAES256Base64Nonce)
}

package main

import (
	"context"

	"github.com/ledongthuc/aws_secrets_storage_sync/cmd"
	"github.com/ledongthuc/aws_secrets_storage_sync/configs"
)

func main() {
	configs.Init()
	configs.PrintConfigs()
	cmd.Start(context.Background())
}

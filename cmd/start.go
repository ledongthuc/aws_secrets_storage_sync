package cmd

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/ledongthuc/aws_secrets_storage_sync/cache"
	"github.com/ledongthuc/aws_secrets_storage_sync/configs"
	"github.com/ledongthuc/aws_secrets_storage_sync/sync"
)

func Start(ctx context.Context) {
	dataSource := cache.NewSecretLastChanges()
	go StartSyncProcess(ctx, dataSource)
	startServer(dataSource)
}

func StartSyncProcess(ctx context.Context, dataSource *cache.SecretLastChanges) {
	syncer := sync.NewSecretSync(dataSource)

	timeoutPeriod := configs.GetSyncPeriodSeconds()
	region := configs.GetAWSRegion()
	if region == "" {
		panic("Missing region config")
	}
	logrus.Infof("Start with %s, duration between runs: %f ms", region, timeoutPeriod)

	filterPrefixName := configs.GetFilterPrefixName()
	filterTags := configs.GetFilterTags()
	filters := sync.BuildAWSFilters(filterPrefixName, filterTags)
	encryption := configs.GetEncryptionConfig()

	for {
		if err := syncer.SyncSecrets(region, filters, filterTags, encryption); err != nil {
			logrus.Warnf("Sync secrets got err: %v", err)
		}
		time.Sleep(time.Duration(timeoutPeriod) * time.Second)
	}
}

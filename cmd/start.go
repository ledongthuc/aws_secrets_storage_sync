package cmd

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/ledongthuc/aws_secrets_storage_sync/configs"
	"github.com/ledongthuc/aws_secrets_storage_sync/sync"
)

func StartSyncProcess(ctx context.Context) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, os.Kill)
	syncer := sync.NewSecretSync()

	timeoutPeriod := configs.GetSyncPeriodSeconds()
	region := configs.GetAWSRegion()
	if region == "" {
		panic("Missing region config")
	}
	logrus.Infof("Start with %s, duration between runs: %f ms", region, timeoutPeriod)

	filterPrefixName := configs.GetFilterPrefixName()
	filterTags := configs.GetFilterTags()
	filters := sync.BuildAWSFilters(filterPrefixName, filterTags)

	for {
		if err := syncer.SyncSecrets(region, filters, filterTags); err != nil {
			logrus.Warnf("Sync secrets got err: %v", err)
		}

		select {
		case <-ctx.Done():
			logrus.Info("Force exit")
			return
		case <-signalChan:
			logrus.Info("Exit")
			return
		case <-time.After(time.Duration(timeoutPeriod) * time.Second):
		}
	}
}

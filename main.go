package main

import (
	"os"
	"os/signal"
	"time"

	"github.com/ledongthuc/aws_secrets_storage_sync/configs"
	"github.com/ledongthuc/aws_secrets_storage_sync/sync"
	"github.com/sirupsen/logrus"
)

func main() {
	configs.Init()
	printConfigs()
	startSyncProcess()
}

func startSyncProcess() {
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
		select {
		case <-signalChan:
			logrus.Info("Exit")
			return
		case <-time.After(time.Duration(timeoutPeriod) * time.Second):
			if err := syncer.SyncSecrets(region, filters, filterTags); err != nil {
				logrus.Warnf("Sync secrets got err: %v", err)
			}
		}
	}
}

func printConfigs() {
	logrus.WithFields(logrus.Fields{
		"AWS_REGION":          configs.GetAWSRegion(),
		"SYNC_PERIOD_SECONDS": configs.GetSyncPeriodSeconds(),
		"FILTER_TAGS":         configs.GetFilterTags(),
		"FILTER_PREFIX_NAME":  configs.GetFilterPrefixName(),
		"SAVING_PATH":         configs.GetSavingPath(),
	}).Info("configs")
}

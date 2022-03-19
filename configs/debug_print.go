package configs

import (
	"github.com/sirupsen/logrus"
)

func PrintConfigs() {
	logrus.WithFields(logrus.Fields{
		"AWS_REGION":          GetAWSRegion(),
		"SYNC_PERIOD_SECONDS": GetSyncPeriodSeconds(),
		"FILTER_TAGS":         GetFilterTags(),
		"FILTER_PREFIX_NAME":  GetFilterPrefixName(),
		"SAVING_PATH":         GetSavingPath(),
	}).Info("configs")
}

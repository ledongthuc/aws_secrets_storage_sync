package configs

import (
	"strings"

	"github.com/spf13/viper"
)

const (
	syncPeriodSeconds = "SYNC_PERIOD_SECONDS"
	filterTags        = "FILTER_TAGS"
	filterPrefixName  = "FILTER_PREFIX_NAME"
	savingPath        = "SAVING_PATH"
)

func GetSyncPeriodSeconds() float64 {
	return viper.GetFloat64(syncPeriodSeconds)
}

func GetFilterTags() [][2]string {
	raw := viper.GetString(filterTags)
	pairs := strings.Split(raw, ",")

	result := make([][2]string, 0, len(pairs))
	for _, pair := range pairs {
		kv := strings.Split(pair, ":")
		if len(kv) != 2 {
			continue
		}
		result = append(result, [2]string{kv[0], kv[1]})
	}
	return result
}

func GetFilterPrefixName() string {
	return viper.GetString(filterPrefixName)
}

func GetSavingPath() string {
	p := viper.GetString(savingPath)
	if !strings.HasSuffix(p, "/") {
		return p + "/"
	}
	return p
}

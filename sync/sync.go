package sync

import (
	"os"
	"time"

	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/ledongthuc/aws_secrets_storage_sync/cache"
	"github.com/ledongthuc/aws_secrets_storage_sync/configs"
	"github.com/ledongthuc/aws_secrets_storage_sync/utils"
)

type SecretSync struct {
	cache *cache.SecretLastChanges
}

func NewSecretSync() *SecretSync {
	return &SecretSync{
		cache: cache.NewSecretLastChanges(),
	}
}

func (s *SecretSync) SyncSecrets(region string, filters []*secretsmanager.Filter, filterTags [][2]string) error {
	secrets, err := GetListSecrets(region, filters, filterTags)
	if err != nil {
		return errors.Wrap(err, "load secrets")
	}

	var total int64
	for _, secret := range secrets {
		if secret == nil || secret.DeletedDate == nil {
			continue
		}
		total++
		// TODO: continue to check sync secret
		if err := s.syncSecret(region, secret); err != nil {
			return errors.Wrapf(err, "sync \"%s\"", utils.Ptr2str(secret.Name))
		}
	}
	logrus.Infof("Sync total %d", total)

	return nil
}

func (s *SecretSync) syncSecret(region string, secret *secretsmanager.SecretListEntry) error {
	if secret == nil {
		return nil
	}

	// check cache and pass if nothing changed
	lastChangeDate, err := getLastChangeDate(secret)
	if err != nil {
		return errors.Wrap(err, "load change date")
	}

	secretName := utils.Ptr2str(secret.Name)
	cachedItem, existed := s.cache.Get(secretName)
	if !existed || !lastChangeDate.After(cachedItem.LastChanged) {
		return nil
	}

	// Saving physical file
	fileName := cachedItem.FileName
	if !existed {
		fileName = utils.RandomString(32)
	}
	savingPath := configs.GetSavingPath()
	err = s.saveSecret(region, secret, savingPath, fileName)
	if err != nil {
		return errors.Wrapf(err, "save file '%s' in path '%s'", fileName, savingPath)
	}

	// Update cache after save successful
	cachedItem.LastChanged = lastChangeDate
	cachedItem.FileName = fileName
	s.cache.Set(secretName, cachedItem)
	return nil
}

func (s *SecretSync) saveSecret(region string, secret *secretsmanager.SecretListEntry, path, fileName string) error {
	// Get secret's value
	value, err := GetSecretValueByARN(region, utils.Ptr2str(secret.ARN))
	if err != nil {
		return errors.Wrap(err, "load secret's value")
	}
	if value == nil {
		return errors.New("value is nil")
	}

	var savingContent []byte
	if value.SecretString != nil {
		savingContent = []byte(*value.SecretString)
	} else if value.SecretBinary != nil {
		savingContent = value.SecretBinary
	} else {
		return errors.New("secret string and binary isn't defined")
	}

	// save to file
	return os.WriteFile(path+fileName, savingContent, 0644)
}

func getLastChangeDate(secret *secretsmanager.SecretListEntry) (time.Time, error) {
	if secret == nil {
		return time.Time{}, errors.New("can't load secret")
	}
	var lastChangeDate time.Time
	if secret.LastChangedDate == nil {
		if secret.CreatedDate == nil {
			return time.Time{}, errors.New("can't get changed date and created date")
		}

		lastChangeDate = *secret.CreatedDate
	} else {
		lastChangeDate = *secret.LastChangedDate
	}
	return lastChangeDate, nil
}

package sync

import (
	"os"
	"path"
	"path/filepath"
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

func NewSecretSync(cache *cache.SecretLastChanges) *SecretSync {
	return &SecretSync{
		cache: cache,
	}
}

func (s *SecretSync) SyncSecrets(region string, filters []*secretsmanager.Filter, filterTags [][2]string, encryption configs.EncryptionConfig) error {
	logrus.Infof("Sync start")
	secrets, err := GetListSecrets(region, filters, filterTags)
	if err != nil {
		return errors.Wrap(err, "load secrets")
	}
	savingPath := configs.GetSavingPath()

	var total int64
	for _, secret := range secrets {
		if secret == nil {
			continue
		}
		err, cached, syncName := s.syncSecret(region, secret, savingPath, encryption)
		if err != nil {
			logrus.Warnf(" - %s: sync failed: %v", utils.Ptr2str(secret.Name), err)
		} else if cached {
			logrus.Infof(" - %s: nothing change: use cache", utils.Ptr2str(secret.Name))
			total++
		} else {
			logrus.Infof(" - %s: sync successful with %s%s", utils.Ptr2str(secret.Name), savingPath, syncName)
			total++
		}
	}
	logrus.Infof("Sync result, total: %d", total)

	return nil
}

func (s *SecretSync) syncSecret(region string, secret *secretsmanager.SecretListEntry, savingPath string, encryption configs.EncryptionConfig) (err error, cached bool, syncName string) {
	if secret == nil {
		return errors.New("error is nil"), false, syncName
	}

	// check cache and pass if nothing changed
	lastChangeDate, err := getLastChangeDate(secret)
	if err != nil {
		return errors.Wrap(err, "load change date"), false, syncName
	}

	secretName := utils.Ptr2str(secret.Name)
	cachedItem, existed := s.cache.Get(secretName)
	if existed && !lastChangeDate.After(cachedItem.LastChanged) {
		return nil, true, syncName
	}

	syncName = cachedItem.FileName
	if !existed {
		syncName = utils.Md5(secretName)
	}

	// Clear cached physical items sync new one
	if existed {
		s.removeOldPhysicalCachedSecret(savingPath, cachedItem.FileName)
	}
	// Saving physical file
	err = s.saveSecret(region, secret, savingPath, syncName, encryption)
	if err != nil {
		return errors.Wrapf(err, "save file '%s' in path '%s'", syncName, savingPath), false, syncName
	}

	// Update cache after save successful
	cachedItem.LastChanged = lastChangeDate
	cachedItem.FileName = syncName
	absLocalPath, err := filepath.Abs(path.Join(savingPath, syncName))
	if err != nil {
		logrus.Warn("can't get abs path %s: %v", path.Join(savingPath, syncName), err)
		cachedItem.LocalPath = path.Join(savingPath, syncName)
	} else {
		cachedItem.LocalPath = absLocalPath
	}
	s.cache.Set(secretName, cachedItem)
	return nil, false, syncName
}

func (s *SecretSync) saveSecret(region string, secret *secretsmanager.SecretListEntry, path, fileName string, encryption configs.EncryptionConfig) error {
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

	if encryption.Method == configs.EncryptionMethodAES256 {
		aes, err := NewAES(encryption.Key, encryption.Nonce)
		if err != nil {
			return errors.Wrap(err, "encrypt saving data")
		}
		savingContent, err = aes.Encrypt(savingContent)
		if err != nil {
			return errors.Wrap(err, "encrypt saving data")
		}
	}

	// save to file
	if err := os.MkdirAll(filepath.Join(path), os.ModePerm); err != nil {
		return err
	}
	return os.WriteFile(path+fileName, savingContent, 0644)
}

func (s *SecretSync) removeOldPhysicalCachedSecret(path, fileName string) error {
	return os.Remove(path + fileName)
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

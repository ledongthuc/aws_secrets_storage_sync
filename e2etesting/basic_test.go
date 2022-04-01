package e2etesting

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"gotest.tools/v3/assert"

	"github.com/ledongthuc/aws_secrets_storage_sync/cmd"
	"github.com/ledongthuc/aws_secrets_storage_sync/configs"
	"github.com/ledongthuc/aws_secrets_storage_sync/utils"
)

func TestBasicStringValue(t *testing.T) {
	os.Setenv("SYNC_PERIOD_SECONDS", "1")
	os.Setenv("AWS_REGION", "ap-southeast-1")
	os.Setenv("ENCRYPTION_METHOD", "NONE")
	os.Setenv("FILTER_PREFIX_NAME", "")
	err := validateEnvironments()
	assert.NilError(t, err, "validate env")

	configs.Init()
	configs.PrintConfigs()

	_, err = prepareAWSSecret(os.Getenv("AWS_REGION"), "string")
	assert.NilError(t, err, "prepare AWS Secret")

	ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)
	cmd.StartSyncProcess(ctx)

	// make sure file are sync to local after 5 seconds
	expectedFileName := utils.Md5(TestSecretName)
	expectedFilePath := fmt.Sprintf("%s/%s", TestSyncPath, expectedFileName)
	if _, err := os.Stat(expectedFilePath); errors.Is(err, os.ErrNotExist) {
		t.Fatalf("%s doesn't existed after sync", expectedFilePath)
	}

	// Check file content
	content, err := os.ReadFile(expectedFilePath)
	assert.NilError(t, err)
	assert.DeepEqual(t, content, []byte(TestSecretStringValue))
}

func TestBasicBinaryValue(t *testing.T) {
	os.Setenv("SYNC_PERIOD_SECONDS", "1")
	os.Setenv("AWS_REGION", "ap-southeast-1")
	os.Setenv("ENCRYPTION_METHOD", "NONE")
	os.Setenv("FILTER_PREFIX_NAME", "")
	err := validateEnvironments()
	assert.NilError(t, err, "validate env")

	configs.Init()
	configs.PrintConfigs()

	_, err = prepareAWSSecret(os.Getenv("AWS_REGION"), "binary")
	assert.NilError(t, err, "prepare AWS Secret")

	ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)
	cmd.StartSyncProcess(ctx)

	// make sure file are sync to local after 5 seconds
	expectedFileName := utils.Md5(TestSecretName)
	expectedFilePath := fmt.Sprintf("%s/%s", TestSyncPath, expectedFileName)
	if _, err := os.Stat(expectedFilePath); errors.Is(err, os.ErrNotExist) {
		t.Fatalf("%s doesn't existed after sync", expectedFilePath)
	}

	// Check file content
	content, err := os.ReadFile(expectedFilePath)
	assert.NilError(t, err)
	assert.DeepEqual(t, content, []byte(TestSecretBinaryValue))
}

func TestFilterName(t *testing.T) {
	os.Setenv("SYNC_PERIOD_SECONDS", "1")
	os.Setenv("AWS_REGION", "ap-southeast-1")
	os.Setenv("ENCRYPTION_METHOD", "NONE")
	os.Setenv("FILTER_PREFIX_NAME", "WRONG_PREFIX")

	err := validateEnvironments()
	assert.NilError(t, err, "validate env")

	configs.Init()
	configs.PrintConfigs()

	_, err = prepareAWSSecret(os.Getenv("AWS_REGION"), "string")
	assert.NilError(t, err, "prepare AWS Secret")

	ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)
	cmd.StartSyncProcess(ctx)

	// File unexisted with wrong prefix filter
	expectedFileName := utils.Md5(TestSecretName)
	expectedFilePath := fmt.Sprintf("%s/%s", TestSyncPath, expectedFileName)
	if _, err := os.Stat(expectedFilePath); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("%s must be not existed with wrong prefix", expectedFilePath)
	}

	os.Setenv("FILTER_PREFIX_NAME", "AWS_SECRET_")
	configs.Init()
	configs.PrintConfigs()

	ctx, _ = context.WithTimeout(context.Background(), 1*time.Second)
	cmd.StartSyncProcess(ctx)

	// make sure file are sync to local after 5 seconds
	expectedFileName = utils.Md5(TestSecretName)
	expectedFilePath = fmt.Sprintf("%s/%s", TestSyncPath, expectedFileName)
	if _, err := os.Stat(expectedFilePath); errors.Is(err, os.ErrNotExist) {
		t.Fatalf("%s doesn't existed after sync", expectedFilePath)
	}

	// Check file content
	content, err := os.ReadFile(expectedFilePath)
	assert.NilError(t, err)
	assert.DeepEqual(t, content, []byte(TestSecretStringValue))
}

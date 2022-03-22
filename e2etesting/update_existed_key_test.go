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

func TestUpdateStringValueOnBinaryValue(t *testing.T) {
	os.Setenv("SYNC_PERIOD_SECONDS", "1")
	err := validateEnvironments()
	assert.NilError(t, err, "validate env")

	configs.Init()
	configs.PrintConfigs()

	_, err = prepareAWSSecret(os.Getenv("AWS_REGION"), "binary")
	assert.NilError(t, err, "prepare AWS Secret")

	go func() {
		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		cmd.StartSyncProcess(ctx)
	}()

	// after 5 seconds, update the secret to string
	time.Sleep(1 * time.Second)
	_, err = prepareAWSSecret(os.Getenv("AWS_REGION"), "string")
	assert.NilError(t, err, "prepare AWS Secret")

	time.Sleep(3 * time.Second)

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

func TestUpdateBinaryValueOnStringValue(t *testing.T) {
	os.Setenv("SYNC_PERIOD_SECONDS", "1")
	err := validateEnvironments()
	assert.NilError(t, err, "validate env")

	configs.Init()
	configs.PrintConfigs()

	_, err = prepareAWSSecret(os.Getenv("AWS_REGION"), "string")
	assert.NilError(t, err, "prepare AWS Secret")

	go func() {
		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		cmd.StartSyncProcess(ctx)
	}()

	// after 5 seconds, update the secret to string
	time.Sleep(1 * time.Second)
	_, err = prepareAWSSecret(os.Getenv("AWS_REGION"), "binary")
	assert.NilError(t, err, "prepare AWS Secret")

	time.Sleep(3 * time.Second)

	// make sure file are sync to local after waiting
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

package e2etesting

import (
	"context"
	"encoding/base64"
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

// TODO: filter test
func TestEncryptStringValue(t *testing.T) {
	os.Setenv("SYNC_PERIOD_SECONDS", "1")
	os.Setenv("AWS_REGION", "ap-southeast-1")
	os.Setenv("ENCRYPTION_METHOD", "AES_GCM_256")
	os.Setenv("FILTER_PREFIX_NAME", "")
	os.Setenv("AES_GCM_256_BASE_64_KEY", "MTIzNDU2Nzg5MDEyMzQ1Njc4OTAxMjM0NTY3ODkwMTI=")
	os.Setenv("AES_GCM_256_BASE_64_NONCE", "MTIzNDU2Nzg5MDEy")

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
	expectedEncodedValue := "9JeOmiI3oI2Qdx9nUhhASMW5ffsa+/RYiAw="
	expectedValue, err := base64.StdEncoding.DecodeString(expectedEncodedValue)
	assert.NilError(t, err)
	assert.DeepEqual(t, content, expectedValue)
}

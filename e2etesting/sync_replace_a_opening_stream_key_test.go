package e2etesting

import (
	"bufio"
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

/*TestUpdateOpeningExistedKey process:
1. Wait to sync process downloads files successfully
2. Test the file + content correctly
3. Open file without closing
4. Update AWS secret directly.
5. Existed opening file has old content (2), new opening file  has new content (4)
*/
func TestUpdateOpeningExistedKey(t *testing.T) {
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

	// 1.
	time.Sleep(2 * time.Second)

	// 2.
	expectedFileName := utils.Md5(TestSecretName)
	expectedFilePath := fmt.Sprintf("%s/%s", TestSyncPath, expectedFileName)
	if _, err := os.Stat(expectedFilePath); errors.Is(err, os.ErrNotExist) {
		t.Fatalf("%s doesn't existed after sync", expectedFilePath)
	}
	content, err := os.ReadFile(expectedFilePath)
	assert.NilError(t, err)
	assert.DeepEqual(t, content, []byte(TestSecretStringValue))

	// 3.
	f1, err := os.Open(expectedFilePath)
	assert.NilError(t, err)
	defer f1.Close()

	// 4.
	_, err = prepareAWSSecret(os.Getenv("AWS_REGION"), "binary")
	assert.NilError(t, err, "prepare AWS Secret")

	time.Sleep(2 * time.Second)

	// 5.
	f2, err := os.Open(expectedFilePath)
	assert.NilError(t, err)
	defer f2.Close()

	scanner := bufio.NewScanner(f1)
	f1Bytes := []byte{}
	for scanner.Scan() {
		f1Bytes = append(f1Bytes, scanner.Bytes()...)
	}
	assert.DeepEqual(t, f1Bytes, []byte(TestSecretStringValue))

	scanner = bufio.NewScanner(f2)
	f2Bytes := []byte{}
	for scanner.Scan() {
		f2Bytes = append(f2Bytes, scanner.Bytes()...)
	}
	assert.DeepEqual(t, f2Bytes, []byte(TestSecretBinaryValue))
}

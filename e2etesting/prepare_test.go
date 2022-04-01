package e2etesting

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var (
	TestSecretName        = "AWS_SECRET_SYNC_E2E_TEST"
	TestSecretTag         = "E2E_TESTING"
	TestSecretStringValue = "test value"
	TestSecretBinaryValue = []byte{1, 2, 3, 4, 5, 6}
	TestSyncPath          = "./file_sync"
)

func prepareAWSSecret(region string, valueType string) (*secretsmanager.SecretListEntry, error) {
	err := os.RemoveAll(TestSyncPath)
	if err != nil {
		return nil, errors.Wrap(err, "remove old files")
	}
	secret, err := getTestAWSSecret(region)
	if err != nil {
		logrus.Infof("Fail to get TEST AWS secret: %v", err)
		secret, err = createTestAWSSecret(region)
		if err != nil {
			return nil, errors.Wrap(err, "create TEST AWS secret")
		}
	}
	if secret.DeletedDate != nil {
		if err := recoverTestAWSSecret(region, secret); err != nil {
			return nil, errors.Wrap(err, "recover TEST AWS secrets")
		}
	}
	updateTestAWSSecret(region, secret, valueType)
	return secret, nil
}

func getTestAWSSecret(region string) (*secretsmanager.SecretListEntry, error) {
	svc := secretsmanager.New(session.New(&aws.Config{
		Region: aws.String(region),
	}))
	input := &secretsmanager.ListSecretsInput{
		MaxResults: aws.Int64(1),
		Filters: []*secretsmanager.Filter{
			&secretsmanager.Filter{
				Key: aws.String("name"),
				Values: []*string{
					aws.String(TestSecretName),
				},
			},
		},
	}
	result, err := svc.ListSecrets(input)
	if err != nil {
		return nil, err
	}
	if len(result.SecretList) == 0 {
		return nil, fmt.Errorf("no secret found")
	}
	return result.SecretList[0], nil
}

func createTestAWSSecret(region string) (*secretsmanager.SecretListEntry, error) {
	svc := secretsmanager.New(session.New(&aws.Config{
		Region: aws.String(region),
	}))
	input := &secretsmanager.CreateSecretInput{
		Name:        aws.String(TestSecretName),
		Description: aws.String("USE FOR E2E TESTING"),
		Tags: []*secretsmanager.Tag{
			{
				Key:   aws.String(TestSecretTag),
				Value: aws.String("true"),
			},
		},
	}
	if _, err := svc.CreateSecret(input); err != nil {
		return nil, err
	}
	return getTestAWSSecret(region)
}

func recoverTestAWSSecret(region string, secret *secretsmanager.SecretListEntry) error {
	svc := secretsmanager.New(session.New(&aws.Config{
		Region: aws.String(region),
	}))
	input := &secretsmanager.RestoreSecretInput{
		SecretId: secret.ARN,
	}
	secret.DeletedDate = nil
	_, err := svc.RestoreSecret(input)
	return err
}

func updateTestAWSSecret(region string, secret *secretsmanager.SecretListEntry, valueType string) error {
	svc := secretsmanager.New(session.New(&aws.Config{
		Region: aws.String(region),
	}))
	input := &secretsmanager.PutSecretValueInput{
		SecretId:     secret.ARN,
		SecretString: aws.String(TestSecretStringValue),
	}
	if valueType == "binary" {
		input.SecretBinary = TestSecretBinaryValue
		input.SecretString = nil
	} else {
		input.SecretBinary = nil
		input.SecretString = aws.String(TestSecretStringValue)
	}
	_, err := svc.PutSecretValue(input)
	return err
}

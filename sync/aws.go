package sync

import (
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/ledongthuc/aws_secrets_storage_sync/utils"
)

func GetListSecrets(region string, filters []*secretsmanager.Filter, filterTags [][2]string) ([]*secretsmanager.SecretListEntry, error) {
	svc := secretsmanager.New(session.New(&aws.Config{
		Region: aws.String(region),
	}))
	fmt.Println("DEBUG: ", region)
	maxResult := int64(100)

	index := 0
	var token *string
	secrets := []*secretsmanager.SecretListEntry{}
	for ; token != nil || index == 0; index++ {
		result, err := GetAPageSecrets(svc, token, maxResult, filters)
		if err != nil {
			return nil, err
		}

		for _, item := range result.SecretList {
			if ValidSecretTags(item, filterTags) {
				secrets = append(secrets, item)
			}
		}
		token = result.NextToken
	}
	return secrets, nil
}

func BuildAWSFilters(filterPrefixName string, filterTags [][2]string) []*secretsmanager.Filter {
	result := []*secretsmanager.Filter{}

	if len(filterPrefixName) > 0 {
		result = append(result, &secretsmanager.Filter{
			Key:    aws.String("name"),
			Values: []*string{aws.String(filterPrefixName)},
		})
	}

	var tagKeys, tagValues []*string
	for _, tag := range filterTags {
		tagKeys = append(tagKeys, &tag[0])
		tagValues = append(tagValues, &tag[1])
	}

	if len(tagKeys) > 0 {
		result = append(result, &secretsmanager.Filter{
			Key:    aws.String("tag-key"),
			Values: tagKeys,
		})
	}

	if len(tagValues) > 0 {
		result = append(result, &secretsmanager.Filter{
			Key:    aws.String("tag-value"),
			Values: tagValues,
		})
	}

	return result
}

func GetAPageSecrets(svc *secretsmanager.SecretsManager, token *string, maxResult int64, filters []*secretsmanager.Filter) (*secretsmanager.ListSecretsOutput, error) {
	input := &secretsmanager.ListSecretsInput{
		MaxResults: &maxResult,
		NextToken:  token,
		Filters:    filters,
	}

	result, err := svc.ListSecrets(input)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func ValidSecretTags(secret *secretsmanager.SecretListEntry, filterTags [][2]string) bool {
	if secret == nil {
		return false
	}
	if len(filterTags) == 0 {
		return true
	}
	for _, tag := range secret.Tags {
		for _, filterTag := range filterTags {
			if utils.Ptr2str(tag.Key) == filterTag[0] && utils.Ptr2str(tag.Value) == filterTag[1] {
				return true
			}
		}
	}
	return false
}

func GetSecretValueByARN(region, arn string) (*secretsmanager.GetSecretValueOutput, error) {
	svc := secretsmanager.New(session.New(&aws.Config{
		Region: aws.String(region),
	}))

	input := &secretsmanager.GetSecretValueInput{SecretId: aws.String(arn)}
	result, err := svc.GetSecretValue(input)

	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, errors.New("Can't get secret")
	}
	return result, nil
}

package utils

import (
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"gopkg.in/ini.v1"
)

var awsCredentialsFile = filepath.Join(os.Getenv("HOME"), ".aws", "credentials")

type awsCredentials struct {
	alias     string
	accessKey string
	secretKey string
	token     string
}

func AssumeAwsRole(roleArn string, sessionName string) error {
	session := sts.New(session.New())

	resp, err := session.AssumeRole(&sts.AssumeRoleInput{
		RoleArn:         &roleArn,
		RoleSessionName: &sessionName,
	})

	if err != nil {
		return err
	}

	saveAwsCredentials(&awsCredentials{
		alias:     sessionName,
		accessKey: *resp.Credentials.AccessKeyId,
		secretKey: *resp.Credentials.SecretAccessKey,
		token:     *resp.Credentials.SessionToken,
	})

	return nil
}

func saveAwsCredentials(awsCredentials *awsCredentials) error {
	cfg := ini.Empty()

	cfg.NewSection(awsCredentials.alias)
	cfg.Section(awsCredentials.alias).NewKey("aws_access_key_id", awsCredentials.accessKey)
	cfg.Section(awsCredentials.alias).NewKey("aws_secret_access_key", awsCredentials.secretKey)
	cfg.Section(awsCredentials.alias).NewKey("aws_security_token", awsCredentials.token)

	err := cfg.SaveTo(awsCredentialsFile)

	if err != nil {
		return err
	}

	return nil
}

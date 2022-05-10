package aws

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"reflect"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/caarlos0/env/v6"
	"github.com/wspowell/context"
	"github.com/wspowell/errors"
	"github.com/wspowell/log"
)

const (
	appName = "snailmail"
	region  = "us-east-1"
)

// GetSecret from either local environment or from AWS SecretsManager.
// AWS secret name is derived from the app name and struct name: <app>-<struct>
//   Example: snailmail-rdsConnectionInfo
func GetSecret(ctx context.Context, secretModel interface{}) error {
	log.Debug(ctx, "getting secret: %T", secretModel)

	if environment, exists := os.LookupEnv("ENV"); exists && environment == "dev" {
		log.Debug(ctx, "dev mode, getting secret from environment")

		// Default to using env var.
		// Useful for local development.
		if err := env.Parse(secretModel); err != nil {
			return err
		}

		return nil
	}

	log.Debug(ctx, "getting secret from SecretsManager")

	secretName := appName + "-" + getTypeName(secretModel)

	log.Debug(ctx, "secretName=%s", secretName)

	// Create a Secrets Manager client
	sess, err := session.NewSession()
	if err != nil {
		// Handle session creation error
		return err
	}
	svc := secretsmanager.New(sess,
		aws.NewConfig().WithRegion(region))
	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secretName),
		VersionStage: aws.String("AWSCURRENT"), // VersionStage defaults to AWSCURRENT if unspecified
	}

	// In this sample we only handle the specific exceptions for the 'GetSecretValue' API.
	// See https://docs.aws.amazon.com/secretsmanager/latest/apireference/API_GetSecretValue.html

	result, err := svc.GetSecretValue(input)
	if err != nil {
		var awsErr awserr.Error
		if errors.As(err, &awsErr) {
			switch awsErr.Code() {
			case secretsmanager.ErrCodeDecryptionFailure:
				// Secrets Manager can't decrypt the protected secret text using the provided KMS key.
				fmt.Println(secretsmanager.ErrCodeDecryptionFailure, awsErr.Error())

			case secretsmanager.ErrCodeInternalServiceError:
				// An error occurred on the server side.
				fmt.Println(secretsmanager.ErrCodeInternalServiceError, awsErr.Error())

			case secretsmanager.ErrCodeInvalidParameterException:
				// You provided an invalid value for a parameter.
				fmt.Println(secretsmanager.ErrCodeInvalidParameterException, awsErr.Error())

			case secretsmanager.ErrCodeInvalidRequestException:
				// You provided a parameter value that is not valid for the current state of the resource.
				fmt.Println(secretsmanager.ErrCodeInvalidRequestException, awsErr.Error())

			case secretsmanager.ErrCodeResourceNotFoundException:
				// We can't find the resource that you asked for.
				fmt.Println(secretsmanager.ErrCodeResourceNotFoundException, awsErr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}

		return err
	}

	// Decrypts secret using the associated KMS CMK.
	// Depending on whether the secret is a string or binary, one of these fields will be populated.
	var secretString string // , decodedBinarySecret string
	if result.SecretString != nil {
		secretString = *result.SecretString
	} else {
		decodedBinarySecretBytes := make([]byte, base64.StdEncoding.DecodedLen(len(result.SecretBinary)))
		length, err := base64.StdEncoding.Decode(decodedBinarySecretBytes, result.SecretBinary)
		if err != nil {
			fmt.Println("Base64 Decode Error:", err)

			return err
		}
		secretString = string(decodedBinarySecretBytes[:length])
	}

	if err := json.Unmarshal([]byte(secretString), &secretModel); err != nil {
		return err
	}

	return nil
}

func getTypeName(value interface{}) string {
	valueType := reflect.TypeOf(value)
	for valueType.Kind() == reflect.Ptr {
		valueType = valueType.Elem()
	}

	return valueType.Name()
}

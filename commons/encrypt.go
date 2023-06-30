package commons

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/dileepaj/tracified-gateway/utilities"
)

//!As the custom logger constants are in the same level please use the direct int value for log levels
//!To break the import cycle

func Encrypt(key string) []byte {
	// Load the Shared AWS Configuration (~/.aws/config)
	logger := utilities.NewCustomLogger()
	svc := kms.New(session.New(&aws.Config{
		Region:      aws.String(GoDotEnvVariable("AWS_REGION")),
		Credentials: credentials.NewStaticCredentials(GoDotEnvVariable("AWS_ACCESS_KEY"), GoDotEnvVariable("AWS_SECRET_KEY"), ""),
	}))

	input := &kms.EncryptInput{
		KeyId:     aws.String(GoDotEnvVariable("AWS_KMS_KEY_ID")),
		Plaintext: []byte(key),
	}

	result, err := svc.Encrypt(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case kms.ErrCodeNotFoundException:
				logger.LogWriter(kms.ErrCodeNotFoundException+aerr.Error(), 3)
			case kms.ErrCodeDisabledException:
				logger.LogWriter(kms.ErrCodeDisabledException+aerr.Error(), 3)
			case kms.ErrCodeKeyUnavailableException:
				logger.LogWriter(kms.ErrCodeKeyUnavailableException+aerr.Error(), 3)
			case kms.ErrCodeDependencyTimeoutException:
				logger.LogWriter(kms.ErrCodeDependencyTimeoutException+aerr.Error(), 3)
			case kms.ErrCodeInvalidKeyUsageException:
				logger.LogWriter(kms.ErrCodeInvalidKeyUsageException+aerr.Error(), 3)
			case kms.ErrCodeInvalidGrantTokenException:
				logger.LogWriter(kms.ErrCodeInvalidGrantTokenException+aerr.Error(), 3)
			case kms.ErrCodeInternalException:
				logger.LogWriter(kms.ErrCodeInternalException+aerr.Error(), 3)
			case kms.ErrCodeInvalidStateException:
				logger.LogWriter(kms.ErrCodeInvalidStateException+aerr.Error(), 3)
			default:
				logger.LogWriter(aerr.Error(), 3)
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			logger.LogWriter("Error when encrypting : "+err.Error(), 3)
		}
		return []byte{}
	}

	//stgOne := string(result.CiphertextBlob)
	// bone := []byte(stgOne)
	// Decrypt(bone)

	return result.CiphertextBlob
}

func Decrypt(arr []byte) string {
	svc := kms.New(session.New(&aws.Config{
		Region:      aws.String(GoDotEnvVariable("AWS_REGION")),
		Credentials: credentials.NewStaticCredentials(GoDotEnvVariable("AWS_ACCESS_KEY"), GoDotEnvVariable("AWS_SECRET_KEY"), ""),
	}))
	logger := utilities.NewCustomLogger()

	input := &kms.DecryptInput{
		CiphertextBlob: arr,
		KeyId:          aws.String(GoDotEnvVariable("AWS_KMS_KEY_ID")),
	}

	result, err := svc.Decrypt(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case kms.ErrCodeNotFoundException:
				logger.LogWriter(kms.ErrCodeNotFoundException+aerr.Error(), 3)
			case kms.ErrCodeDisabledException:
				logger.LogWriter(kms.ErrCodeDisabledException+aerr.Error(), 3)
			case kms.ErrCodeInvalidCiphertextException:
				logger.LogWriter(kms.ErrCodeInvalidCiphertextException+aerr.Error(), 3)
			case kms.ErrCodeKeyUnavailableException:
				logger.LogWriter(kms.ErrCodeKeyUnavailableException+aerr.Error(), 3)
			case kms.ErrCodeIncorrectKeyException:
				logger.LogWriter(kms.ErrCodeIncorrectKeyException+aerr.Error(), 3)
			case kms.ErrCodeInvalidKeyUsageException:
				logger.LogWriter(kms.ErrCodeInvalidKeyUsageException+aerr.Error(), 3)
			case kms.ErrCodeDependencyTimeoutException:
				logger.LogWriter(kms.ErrCodeDependencyTimeoutException+aerr.Error(), 3)
			case kms.ErrCodeInvalidGrantTokenException:
				logger.LogWriter(kms.ErrCodeInvalidGrantTokenException+aerr.Error(), 3)
			case kms.ErrCodeInternalException:
				logger.LogWriter(kms.ErrCodeInternalException+aerr.Error(), 3)
			case kms.ErrCodeInvalidStateException:
				logger.LogWriter(kms.ErrCodeInvalidStateException+aerr.Error(), 3)
			default:
				logger.LogWriter("Error when decrypting : "+aerr.Error(), 3)
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			logger.LogWriter("Error when decrypting : "+err.Error(), 3)
		}
		return ""
	}
	return string(result.Plaintext)
}

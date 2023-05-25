package ipfsservices

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/dileepaj/tracified-gateway/commons"
)

type IPFSService struct {
}

/*
Create a new bucket - Add the new bucket name to the IPFS Enum file as well before deploying
*/
func (s *IPFSService) CreateBucket(bucketName string) error {
	accessKey := commons.GoDotEnvVariable("FILEBASE_ACCESS_KEY")
	secretKey := commons.GoDotEnvVariable("FILEBASE_SECRET_KEY")
	endpoint := commons.GoDotEnvVariable("FILEBASE_S3_API_ENDPOINT")
	region := commons.GoDotEnvVariable("FILEBASE_REGION")
	profile := commons.GoDotEnvVariable("FILEBASE_PROFILE")

	//create the configuration
	s3Config := aws.Config{
		Credentials:      credentials.NewStaticCredentials(accessKey, secretKey, ""),
		Endpoint:         aws.String(endpoint),
		Region:           aws.String(region),
		S3ForcePathStyle: aws.Bool(true),
	}

	goSession, errWhenCreatingSession := session.NewSessionWithOptions(session.Options{
		Config:  s3Config,
		Profile: profile,
	})
	if errWhenCreatingSession != nil {
		fmt.Println("Error when creating session : ", errWhenCreatingSession.Error())
		return errWhenCreatingSession
	}

	//create the S3 client session
	s3Client := s3.New(goSession)

	//set parameter for the bucket name
	bucket := aws.String(bucketName)

	//create the bucket
	_, errWhenCreatingBucket := s3Client.CreateBucket(&s3.CreateBucketInput{
		Bucket: bucket,
	})
	if errWhenCreatingBucket != nil {
		fmt.Println("Error when creating bucket : ", errWhenCreatingBucket.Error())
		return errWhenCreatingBucket
	}
	return nil
}

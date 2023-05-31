package ipfsservices

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/dileepaj/tracified-gateway/commons"
)

//Upload the content to filebase IPFS bucket
//
//Parameters - File body, Key name, Bucket Name
func (s *IPFSService) UploadFile(fileBody []byte, keyName string, bucketName string) (string, error) {
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
		return "", errWhenCreatingSession
	}

	//create the S3 client session
	s3Client := s3.New(goSession)

	// create put object input
	putObjectInput := &s3.PutObjectInput{
		Body:   bytes.NewReader(fileBody),
		Bucket: aws.String(bucketName),
		Key:    aws.String(keyName),
	}

	// upload file
	_, errWhenUploadingTheFile := s3Client.PutObject(putObjectInput)
	if errWhenUploadingTheFile != nil {
		fmt.Println("Error when uploading the file to the ", bucketName, " bucket in filebase : ", errWhenUploadingTheFile.Error())
		return "", errWhenUploadingTheFile
	}

	resp, errWhenGettingHeadObject := s3Client.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(keyName),
	})
	if errWhenGettingHeadObject != nil {
		fmt.Println("Error when getting the header object : ", errWhenGettingHeadObject)
		return "", errWhenGettingHeadObject
	}

	cid := ""
	if resp.Metadata != nil {
		cidValue, ok := resp.Metadata["Cid"]
		if !ok {
			fmt.Println("Error when getting CID ")
		}
		cid = *cidValue

	} else {
		fmt.Println("CID is not created")
		return "", errors.New("No CID is created")
	}

	fmt.Println("Content uploaded to IPFS at : https://ipfs.filebase.io/ipfs/" + cid)
	return cid, nil
}

package serviceimplementation

import (
	"bytes"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	log "github.com/sirupsen/logrus"
)

// TranscribeJobResult contains the transcription results after downloading. Definition incompleted, only minimum
type TranscribeJobResult struct {
	Results TranscribeResults `json:"results"`
}

// TranscribeResults contain actually the transcript as an array. Definition incompleted, only minimum
type TranscribeResults struct {
	Transcripts []Transcript `json:"transcripts"`
}

// Transcript contain actually the transcript as one string. Definition incompleted, only minimum
type Transcript struct {
	Transcript string `json:"transcript"`
}

func uploadToS3(sess *session.Session, input []byte, myKey, filename string) error {
	svc := s3.New(sess)
	// first check if there is already such a item in the bucket
	loi := &s3.ListObjectsV2Input{
		Bucket: aws.String("dl.vassiliou-pohl.berlin/"),
	}
	objectList, err := svc.ListObjectsV2(loi)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchBucket:
				return aerr
			default:
				return aerr
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			return err
		}
	}
	log.Tracef("Iterating over BucketObjects: %#v", objectList)
	for i, eachObj := range objectList.Contents {
		log.Tracef("BucketElement[%d]: %s\n", i, *eachObj.Key)
		if *eachObj.Key == ("transcribe/" + myKey) {
			log.Tracef("Element found %s. No need to upload\n", *eachObj.Key)
			return nil
		}
	}
	// Create an uploader with the session and default options
	uploader := s3manager.NewUploader(sess)

	// Upload the file to S3.
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String("dl.vassiliou-pohl.berlin/transcribe/"),
		Key:    aws.String(myKey),
		Body:   bytes.NewReader(input),
	})
	if err != nil {
		return fmt.Errorf("failed to upload file, %v", err)
	}
	fmt.Printf("file uploaded to, %s\n", result.Location)
	return nil
}

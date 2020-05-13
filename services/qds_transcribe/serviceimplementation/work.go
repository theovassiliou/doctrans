package serviceimplementation

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/ArthurHlt/go-eureka-client/eureka"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	ts "github.com/aws/aws-sdk-go/service/transcribeservice"
	log "github.com/sirupsen/logrus"
	pb "github.com/theovassiliou/doctrans/dtaservice"
)

type DtaService struct {
	pb.UnimplementedDTAServerServer
	pb.DocTransServer
	resolver              *eureka.Client
	listener              net.Listener
	AWSRegion             string
	AWSCredentialFileName string
	AWSCredentialProfile  string
	AWSBucketPath         string
}

func check(e error) {
	if e != nil {
		log.Errorln(e)
	}
}

var sess *session.Session

// Work uses AWS service to transcribe an audio
func Work(s *DtaService, input []byte, fileName string) (string, []string, error) {
	rand.Seed(time.Now().UTC().UnixNano())
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("eu-central-1"),
		Credentials: credentials.NewSharedCredentials("", "transscribe-test-account"),
	})
	if err != nil {
		log.WithFields(log.Fields{"Service": s.ApplicationName(), "Task": "TransformDocument.Work"}).Errorln("Can not create AWS session. " + err.Error())
	}

	// Check whether file is supported file format
	for err := checkFileType(input, fileName); err != nil; {
		// No --> Return error
		log.WithFields(log.Fields{"Service": s.ApplicationName(), "Task": "TransformDocument.Work"}).Errorln(err.Error())
		return "", nil, err
	}
	theJob := &ts.StartTranscriptionJobInput{
		LanguageCode: aws.String(ts.LanguageCodeDeDe),
		MediaFormat:  aws.String(ts.MediaFormatMp3),
	}
	ressourceItemName := calcRessourceName(input, fileName, theJob)

	theJob = &ts.StartTranscriptionJobInput{
		LanguageCode: aws.String(ts.LanguageCodeDeDe),
		Media: &ts.Media{
			MediaFileUri: aws.String("s3://dl.vassiliou-pohl.berlin/transcribe/" + ressourceItemName),
		},
		MediaFormat: aws.String(ts.MediaFormatMp3),
	}

	// Save file in S3 bucket
	for err := uploadToS3(sess, input, ressourceItemName, fileName); err != nil; {
		log.WithFields(log.Fields{"Service": s.ApplicationName(), "Task": "TransformDocument.Work"}).Errorln(err.Error())
		return "", nil, err
	}

	// Start Transcription Job
	// Return result
	svc := ts.New(sess)

	jobName := ressourceItemName
	theJob.SetTranscriptionJobName(jobName)

	_, err = svc.StartTranscriptionJob(theJob)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ts.ErrCodeConflictException:
				log.Traceln("Transcribing with jobName:", jobName)
				log.Debugln("Job already exists. No need to redo. Reusing results")
				break
			default:
				log.Traceln("Transcribing with jobName:", jobName)
				log.Debugf("Error occured on starting Transcription. %s", err.Error())
				return "", nil, err
			}
		}
	}

	j, _ := svc.GetTranscriptionJob(&ts.GetTranscriptionJobInput{
		TranscriptionJobName: theJob.TranscriptionJobName,
	})

	for *j.TranscriptionJob.TranscriptionJobStatus != ts.TranscriptionJobStatusCompleted && *j.TranscriptionJob.TranscriptionJobStatus != ts.TranscriptionJobStatusFailed {
		time.Sleep(5 * time.Second)
		j, _ = svc.GetTranscriptionJob(&ts.GetTranscriptionJobInput{
			TranscriptionJobName: theJob.TranscriptionJobName,
		})
		log.Tracef("TranscriptJob: %v", j)
	}

	if *j.TranscriptionJob.TranscriptionJobStatus == ts.TranscriptionJobStatusFailed {
		return "", []string{}, fmt.Errorf("Failed to transform document. %s", *j.TranscriptionJob.FailureReason)
	}

	if j.TranscriptionJob.Transcript.TranscriptFileUri != nil {

		resp, err := http.Get(*j.TranscriptionJob.Transcript.TranscriptFileUri)
		if err != nil {
			// handle error
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)

		var f TranscribeJobResult
		json.Unmarshal(body, &f)
		if f.Results.Transcripts != nil {
			var b strings.Builder
			for _, trans := range f.Results.Transcripts {
				b.WriteString(trans.Transcript)
			}
			return b.String(), []string{string(body)}, nil
		}
		return string(body), []string{}, nil
	}
	return "", nil, nil

}
func ExecuteWorkerLocally(s DtaService, fileName string, additionalInfo bool) {
	if fileName == "" {
		log.Errorln("No fileName on local executing provided. Aborting.")
		return
	}

	dat, err := ioutil.ReadFile(fileName)
	check(err)

	transDoc, stdout, err := Work(&s, dat, fileName)
	if err != nil {
		log.Errorln(err.Error())
		os.Exit(1)
	}

	fmt.Printf("Transcription Results:\n%s\n", transDoc)

	if additionalInfo {
		for i, aResult := range stdout {
			fmt.Printf("out[%d]: %s\n", i, aResult)
		}
	}
}

func calcRessourceName(input []byte, fileName string, theJob *ts.StartTranscriptionJobInput) string {
	if input == nil {
		return fileName
	}
	h := hmac.New(sha256.New, []byte("P4HG#BjA3S85"))
	h.Write(input)
	b, _ := json.Marshal(*theJob)
	h.Write(b)
	return hex.EncodeToString(h.Sum(nil))
}

func checkFileType(input []byte, fileName string) error {
	if fileName != "" {
		if strings.HasSuffix(fileName, "mp3") || strings.HasSuffix(fileName, "wav") {
			return nil
		}
		return errors.New("File suffix not supported. Neither mp3 nor wav")
	}
	return nil
}

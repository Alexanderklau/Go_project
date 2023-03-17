package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}

func main() {

	access_key := "JCO64M9W5TLCTDKD35O1"
	secret_key := "hPmZaZpxLqHyBJDPN3fUAskZhOUT71xWB3eyW6rX"
	end_point := "http://10.0.6.247:7480" //endpoint设置，不要动

	sess, err := session.NewSession(&aws.Config{
		Credentials:      credentials.NewStaticCredentials(access_key, secret_key, ""),
		Endpoint:         aws.String(end_point),
		Region:           aws.String("us-east-1"),
		DisableSSL:       aws.Bool(true),
		S3ForcePathStyle: aws.Bool(false), //virtual-host style方式，不要修改
	})

	bucket := os.Args[1]

	if len(os.Args) != 2 {
		exitErrorf("Bucket name required\nUsage: %s bucket_name",
			os.Args[0])
	}

	// Create S3 service client
	svc := s3.New(sess)

	params := &s3.CreateBucketInput{
		Bucket: aws.String(bucket),
	}

	_, err = svc.CreateBucket(params)

	if err != nil {
		exitErrorf("Unable to create bucket %q, %v", bucket, err)
	}

	// Wait until bucket is created before finishing
	fmt.Printf("Waiting for bucket %q to be created...\n", bucket)

	err = svc.WaitUntilBucketExists(&s3.HeadBucketInput{
		Bucket: aws.String(bucket),
	})

	if err != nil {
		exitErrorf("Error occurred while waiting for bucket to be created, %v", bucket)
	}

	fmt.Printf("Bucket %q successfully created\n", bucket)
}

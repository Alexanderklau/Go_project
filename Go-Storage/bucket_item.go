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
		Region:           aws.String("us-west-1"),
		DisableSSL:       aws.Bool(true),
		S3ForcePathStyle: aws.Bool(false), //virtual-host style方式，不要修改
	})

	if len(os.Args) != 2 {
		exitErrorf("Bucket name required\nUsage: %s bucket_name",
			os.Args[0])
	}

	bucket := os.Args[1]
	fmt.Printf(bucket)
	fmt.Printf("\n")

	// Create S3 service client
	svc := s3.New(sess)

	params := &s3.ListObjectsInput{
		Bucket: aws.String(bucket),
	}
	resp, err := svc.ListObjects(params)

	if err != nil {
		exitErrorf("Unable to list items in bucket %q, %v", bucket, err)
	}

	for _, item := range resp.Contents {
		fmt.Println("Name:         ", *item.Key)
		fmt.Println("Last modified:", *item.LastModified)
		fmt.Println("Size:         ", *item.Size)
		fmt.Println("Storage class:", *item.StorageClass)
		fmt.Println("")
	}
}

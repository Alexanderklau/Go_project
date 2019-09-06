package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
    "fmt"
    "os"
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
        Region:           aws.String("us-west-2"),
        DisableSSL:       aws.Bool(true),
		S3ForcePathStyle: aws.Bool(false), //virtual-host style方式，不要修改
	})

	if len(os.Args) != 3 {
		exitErrorf("bucket and file name required\nUsage: %s bucket_name filename",
			os.Args[0])
	}
	
	bucket := os.Args[1]
	filename := os.Args[2]
	
	file, err := os.Open(filename)
	if err != nil {
		exitErrorf("Unable to open file %q, %v", err)
	}
	
	defer file.Close()

	uploader := s3manager.NewUploader(sess)

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key: aws.String(filename),
		Body: file,
	})
	if err != nil {
		// Print the error and exit.
		exitErrorf("Unable to upload %q to %q, %v", filename, bucket, err)
	}
	
	fmt.Printf("Successfully uploaded %q to %q\n", filename, bucket)
}
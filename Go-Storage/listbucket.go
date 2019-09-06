package main
import (
	"fmt"
	"os"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/credentials"
    _ "github.com/aws/aws-sdk-go/service/s3/s3manager"
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
	svc := s3.New(sess)
	result, err := svc.ListBuckets(nil)
	if err != nil {
		exitErrorf("Unable to list buckets, %v", err)
	}
	
	fmt.Println("Buckets:")
	
	for _, b := range result.Buckets {
		fmt.Printf("* %s created on %s\n",
			aws.StringValue(b.Name), aws.TimeValue(b.CreationDate))
	}

	for _, b := range result.Buckets {
		fmt.Printf("%s\n", aws.StringValue(b.Name))
	}
	
}
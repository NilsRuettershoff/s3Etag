package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/CologneBroadcastingCenter/s3Etag"
)

func main() {
	var path, bucket, key, region string
	flag.StringVar(&path, "path", "", "path to file for checksum")
	flag.StringVar(&bucket, "bucket", "", "bucket of object")
	flag.StringVar(&key, "key", "", "key of object")
	flag.StringVar(&region, "region", "", "region of object")
	flag.Parse()
	fsEtag := ""
	if path != "" {
		tag, err := s3Etag.CalculateLocalETag(path, 5)
		if err != nil {
			log.Printf("err: %v\n", err)
			os.Exit(1)
		}
		fsEtag = tag
	}
	awsEtag := ""
	if bucket != "" && key != "" && region != "" {
		session := session.Must(session.NewSession(&aws.Config{
			Region: aws.String(region),
		}))
		s3client := s3.New(session)
		tag, err := s3Etag.FetchS3Etag(s3client, bucket, key)
		if err != nil {
			log.Printf("err: %v\n", err)
			os.Exit(1)
		}
		awsEtag = tag
	}
	if awsEtag != "" && fsEtag != "" {
		comp := "differ"
		if awsEtag == fsEtag {
			comp = "same"
		}
		fmt.Printf("FS Etag: %s | S3 Etag: %s | comparission: %s\n", fsEtag, awsEtag, comp)
	}
	if awsEtag != "" && fsEtag == "" {
		fmt.Printf("S3 Etag: %s\n", awsEtag)
	}

	if fsEtag != "" && awsEtag == "" {
		fmt.Printf("FS Etag: %s\n", fsEtag)
	}
}

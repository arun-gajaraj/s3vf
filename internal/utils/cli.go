package utils

import (
	"fmt"
	"os"
)

func SetArgs(s3c *S3Config) {
	notSet := os.Getenv(`S3_BUCKET_NAME`) == `` ||
		os.Getenv(`S3_BUCKET_REGION`) == `` ||
		os.Getenv(`S3_OBJECT_KEY`) == ``

	if !notSet {
		s3c.Bucket = os.Getenv(`S3_BUCKET_NAME`)
		s3c.Region = os.Getenv(`S3_BUCKET_REGION`)
		s3c.Key = os.Getenv(`S3_OBJECT_KEY`)
		return
	}

	if args := os.Args[1:]; len(args) == 3 {

		s3c.Bucket = args[0]
		s3c.Region = args[1]
		s3c.Key = args[2]

	} else {
		var b, r, k string
		fmt.Println("Enter the s3 bucket name:")
		fmt.Scanln(&b)
		fmt.Println("Enter the s3 bucket region:")
		fmt.Scanln(&r)
		fmt.Println("Enter the s3 object name:")
		fmt.Scanln(&k)

		s3c.Bucket = b
		s3c.Region = r
		s3c.Key = k
	}
}

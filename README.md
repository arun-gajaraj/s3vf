# s3vf
A simple command line utility to download the old versions of a file from s3 when the aws console can't help.


#### Pre-req:

Env Vars: 
- AWS_ACCESS_KEY_ID
- AWS_SECRET_ACCESS_KEY
- AWS_SESSION_TOKEN
- S3_BUCKET_NAME    [optional]
- S3_BUCKET_REGION  [optional]
- S3_OBJECT_KEY     [optional]
- INDENT_JSON       [optional; example: true]<br>

#### to run from source:

`go run s3vf.go <aws_s3_bucket_name> <aws_s3_bucket_region> <s3_filename>`

or <br>

`go run s3vf.go`<br>

### How it works?

1. Set the Variables and run <br>
2. Enter the Date and Time Range and the versions within that time will be downloaded to ./downloads <br>
  <br>
 - Files will be named: [version-last-modified-time] - [version-id].extension <br>
 - Uses Local Time <br>
 - Option to have the JSON file Indented <br>

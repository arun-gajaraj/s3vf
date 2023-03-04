# s3vf
A simple command line utility to download the versions of an s3 object versions with their Last Modified Date.


#### Pre-req:

Set the env vars for aws: 
- AWS_ACCESS_KEY_ID
- AWS_SECRET_ACCESS_KEY<br>

#### to run from source:

`go run s3vf.go <aws_s3_bucket_name> <aws_s3_bucket_region> <s3_filename>`

or 

#### To run from binary:

`s3vf <aws_s3_bucket_name> <aws_s3_bucket_region> <s3_filename>`

or 

#### run w/o args to enter manually

`go run s3vf.go`<br>
or <br>
`s3vf`
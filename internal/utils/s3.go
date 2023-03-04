package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/adhocore/chin"
	"github.com/arun-gajaraj/s3vf/internal/constants"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	log "github.com/sirupsen/logrus"
)

var wg sync.WaitGroup

type S3Config struct {
	Bucket string
	Region string
	Key    string
}

func GetAllVersionsTill(s3c *S3Config, d time.Time) []*s3.ObjectVersion {
	spin := chin.New()
	go spin.Start()
	defer spin.Stop()

	var allVersions []*s3.ObjectVersion

	// Initialize a new session in the us-west-2 region
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(s3c.Region),
	}))

	// Create a new S3 client
	svc := s3.New(sess)

	// Prepare the input parameters for ListObjectVersions
	input := &s3.ListObjectVersionsInput{
		Bucket: aws.String(s3c.Bucket),
		Prefix: aws.String(s3c.Key),
	}
	fmt.Println(`Getting the versions metadata:`)

	// List the versions of the object from S3
	err := svc.ListObjectVersionsPages(input, func(page *s3.ListObjectVersionsOutput, lastPage bool) bool {
		allVersions = append(allVersions, page.Versions...)

		if page.Versions[0].LastModified.Local().Before(d.Local()) {

			log.Info(`got enough versions till ` + page.Versions[0].LastModified.Local().String())
			return false
		}

		// CAN REMOVE THIS
		fmt.Println(len(allVersions))
		if len(allVersions) > constants.MaxVersions {
			log.Errorf("version too old, crossing %d versions earlier, stopping", constants.MaxVersions)
			return false
		}

		input.KeyMarker = page.NextKeyMarker
		input.VersionIdMarker = page.NextVersionIdMarker
		return !lastPage
	})
	if err != nil {
		panic(err)
	}

	return allVersions
}

func DownloadVersions(cfg *S3Config, versionsToDl []*s3.ObjectVersion) {
	log.Info(fmt.Sprintf(`Downloads Started for %d files`, len(versionsToDl)))

	var workerCount int

	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(cfg.Region),
	}))

	svc := s3.New(sess)

	if len(versionsToDl) < 100 {
		workerCount = 5
	} else {
		workerCount = 25
	}

	jobs := make(chan *s3.ObjectVersion, 50)
	results := make(chan *s3.ObjectVersion, len(versionsToDl))

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go downloader(svc, cfg, jobs, results)
	}

	for _, v := range versionsToDl {
		jobs <- v
	}
	close(jobs)
	wg.Wait()
	log.Info("Files Downloaded")
}

func downloader(svc *s3.S3, cfg *S3Config, jobs <-chan *s3.ObjectVersion, results chan<- *s3.ObjectVersion) {
	defer wg.Done()
	for job := range jobs {

		input := &s3.GetObjectInput{
			Bucket:    aws.String(cfg.Bucket),
			Key:       aws.String(*job.Key),
			VersionId: aws.String(*job.VersionId),
		}

		result, err := svc.GetObject(input)
		if err != nil {
			fmt.Println("error Getting object:", err)
		}

		cwd, _ := os.Getwd()
		err = os.MkdirAll("downloads", os.ModePerm)
		if err != nil {
			log.WithError(err).Error(`error in creating downloads directory`)
		}
		filename := job.LastModified.Local().String() + " - " + *job.VersionId + filepath.Ext(cfg.Key)
		filename = strings.ReplaceAll(filename, ":", "-")

		file, err := os.Create(cwd + string(os.PathSeparator) + "downloads" + string(os.PathSeparator) + filename)
		if err != nil {
			fmt.Println("error creating file:", err)
		}

		indent, err := strconv.ParseBool(os.Getenv(`INDENT_JSON`))
		if err != nil {
			log.WithError(err).Error(`error parsing env var INDENT_JSON`)
		}

		if constants.IndentJSON || indent {
			var encodedStr map[string]interface{}
			encoder := json.NewEncoder(file)
			encoder.SetIndent("", "    ")
			str, err := io.ReadAll(result.Body)
			if err != nil {
				log.WithError(err).Error(`error reading result body`)
			}
			err = json.Unmarshal(str, &encodedStr)
			if err != nil {
				log.WithError(err).Error(`error unmarshalling json`)
			}
			err = encoder.Encode(encodedStr)
			if err != nil {
				log.WithError(err).Error(`error encoding json`)
			}
		} else {
			_, err = io.Copy(file, result.Body)
			if err != nil {
				fmt.Println("error writing to file:", err)
			}
		}
		file.Close()
		result.Body.Close()

	}
}

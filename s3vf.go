package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/adhocore/chin"
	"github.com/arun-gajaraj/s3vf/internal/utils"
	"github.com/aws/aws-sdk-go/service/s3"
	log "github.com/sirupsen/logrus"
)

var (
	date time.Time
	s3c  = &utils.S3Config{}
)

var versions, versionsToDownload []*s3.ObjectVersion

func main() {
	log.Info("Hello!")
	utils.SetArgs(s3c)

	if os.Getenv("AWS_ACCESS_KEY_ID") == "" || os.Getenv("AWS_SECRET_ACCESS_KEY") == "" {
		log.Fatal("AWS_ACCESS_KEY_ID/AWS_SECRET_ACCESS_KEY/AWS_SESSION_TOKEN is not set. set these variables and try again")
	}

	GetInpDate(&date)
	from, to := GetInpTime(date)
	versions = utils.GetAllVersionsTill(s3c, date)

	spin := chin.New()
	go spin.Start()
	defer spin.Stop()

	versionsToDownload = filterVersions(versions, from, to)
	utils.DownloadVersions(s3c, versionsToDownload)
}

func filterVersions(allVersions []*s3.ObjectVersion, from time.Time, to time.Time) []*s3.ObjectVersion {
	var versionsToDownload []*s3.ObjectVersion

	for _, version := range allVersions {
		if !version.LastModified.Local().Before(from) && !version.LastModified.Local().After(to) {
			versionsToDownload = append(versionsToDownload, version)
		}
	}
	return versionsToDownload
}

func GetInpDate(dateToSet *time.Time) string {
	var dateStr string
	fmt.Println("Enter date of the version(s) to get [YYYY-MM-DD]:")
	fmt.Scanln(&dateStr)

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		log.WithError(err).Error("error parsing date")
	}
	*dateToSet = date
	log.Info("Date entered: ", date.Format("2006-01-02"))
	return dateStr
}

func GetInpTime(date time.Time) (time.Time, time.Time) {
	var timeStr string
	var fromTime, toTime time.Time
	locationTime, _ := time.LoadLocation("Local")

	fmt.Println(`Enter Time/Time Range to download [ex: 1530 or 1530-1545]: `)
	fmt.Scanln(&timeStr)

	if split := strings.Split(timeStr, "-"); len(split) == 1 {
		fromTime, _ = time.Parse("1504", split[0])
		toTime = fromTime.Add(time.Minute)
	} else if len(split) == 2 {
		fromTime, _ = time.Parse("1504", split[0])
		toTime, _ = time.Parse("1504", split[1])
	}

	return time.Date(date.Year(), date.Month(), date.Day(), fromTime.Hour(), fromTime.Minute(), fromTime.Second(), 0, locationTime),
		time.Date(date.Year(), date.Month(), date.Day(), toTime.Hour(), toTime.Minute(), toTime.Second(), 0, locationTime)
}

package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type local struct {
	region string
	bucket string
}

func setLocals(locals *local) {
	locals.region = "bar"          // set here your region
	locals.bucket = "foo" // set here your bucket name
}

func download(sess *session.Session, bucket, key string) {
	downloader := s3manager.NewDownloader(sess)

	localFilePath := filepath.Join("files", key)

	if err := os.MkdirAll(filepath.Dir(localFilePath), os.ModePerm); err != nil {
		log.Println("Failed to create dir", err)
		return
	}

	file, err := os.Create(localFilePath)
	if err != nil {
		log.Println("Failed to create file", err)
		return
	}

	defer file.Close()

	numBytes, err := downloader.Download(file, &s3.GetObjectInput{
		Bucket: aws.String(bucket), Key: aws.String(key),
	})

	if err != nil {
		log.Println("Failed to download file", err)
		return
	}

	log.Println("[OK] Downloaded: ", file.Name(), " bytes: ", numBytes)
}

func registerNotFound(files []string) {
	file, _ := os.Create("not_found.txt")
	defer file.Close()

	for _, name := range files{
		file.WriteString(name + "\n")
	}
}

func compareFiles() {
	pathXVM := "xvm"
	pathFiles := "files"
	var xvmFileNames []string
	var downloadedFileNames []string

	filepath.Walk(pathXVM, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("Failed to read directory %s: %v\n", pathXVM, err)
			return err
		}
		if !info.IsDir() {
			xvmFileNames = append(xvmFileNames, info.Name())
		}
		return nil
	})
	//log.Println("XVMS ", xvmFileNames)
	filepath.Walk(pathFiles, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("Failed to read directory %s: %v\n", pathFiles, err)
			return err
		}
		if !info.IsDir() {
			downloadedFileNames = append(downloadedFileNames, info.Name())
		}
		return nil
	})

	var notFound []string

	//log.Println("DOWNLOADES ", downloadedFileNames)
	// file exits in bucket but not in git
	for _, fileBucket := range downloadedFileNames {
		found := false
		for _, fileGit := range xvmFileNames {
			if strings.Split(fileGit, ".")[0] == strings.Split(fileBucket, ".")[0] {
				found = true
			}
		}
		if !found {
			notFound = append(notFound, fileBucket)
		}
	}

	registerNotFound(notFound)
}

func main() {
	// log.Println("Starting...")

	// locales := &local{}
	// setLocals(locales)

	// log.Println("Setting region: ", locales.region, " on bucket ", locales.bucket)

	// _session, err := session.NewSession(&aws.Config{
	// 	Region: aws.String(locales.region),
	// })

	// if err != nil {
	// 	log.Fatal("Failed to create a session ", err)
	// 	return
	// }

	// log.Println("Creating a new s3 client")
	// s3Client := s3.New(_session)
	// // list s3 files
	// err = s3Client.ListObjectsV2Pages(&s3.ListObjectsV2Input{
	// 	Bucket: aws.String(locales.bucket),
	// }, func(page *s3.ListObjectsV2Output, lastPage bool) bool {
	// 	for _, item := range page.Contents {
	// 		download(_session, locales.bucket, *item.Key)
	// 	}
	// 	return !lastPage
	// })
	// if err != nil {
	// 	log.Fatal("Fail to list objects", err)
	// 	return
	// }
	compareFiles()

}

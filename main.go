package main

import (
	"fmt"
	"log"
	"os"
	//
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "s3 plugin"
	app.Usage = "s3 plugin"
	flags := []cli.Flag{
		cli.StringFlag{
			Name:   "bucket",
			Usage:  "bucket for the s3 object",
			EnvVar: "S3_BUCkET",
		},
		cli.StringFlag{
			Name:  "file-name",
			Usage: "source file to send",
		},
		cli.StringFlag{
			Name:   "region",
			Usage:  "aws region",
			Value:  "us-east-1",
			EnvVar: "AWS_DEFAULT_REGION,AWS_REGION",
		},
		cli.StringFlag{
			Name:   "prefix",
			Usage:  "upload files to target folder",
			EnvVar: "S3_PREFIX",
		},
	}
	app.Version = fmt.Sprintf("1.0")

	app.Commands = []cli.Command{
		{
			Name:    "put",
			Aliases: []string{"p"},
			Usage:   "copy a file to s3",
			Action:  s3put,
			Flags:   flags,
		},
		{
			Name:    "get",
			Aliases: []string{"g"},
			Usage:   "get a file from s3",
			Action:  s3get,
			Flags:   flags,
		},
	}
	app.Run(os.Args)
}
func s3put(c *cli.Context) error {
	sess := session.Must(session.NewSession())
	// Create an uploader with the session and default options
	uploader := s3manager.NewUploader(sess)

	f, err := os.Open(c.String("file-name"))
	if err != nil {
		return cli.NewExitError("failed to open file '"+c.String("file-name")+"'", 86)
	}

	if len(c.String("bucket")) == 0 {
		return cli.NewExitError("No bucket specified.", 87)
	}
	if len(c.String("prefix")) == 0 {
		return cli.NewExitError("No prefix specified.", 87)
	}

	// Upload the file to S3.
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(c.String("bucket")),
		Key:    aws.String(c.String("prefix") + "/" + c.String("file-name")),
		Body:   f,
	})

	if err != nil {
		return cli.NewExitError("Failed to upload the file.\n", 87)
	}
	fmt.Printf("File succesffully uploaded to %s\n", result.Location)
	return nil
}

func s3get(c *cli.Context) error {
	sess := session.Must(session.NewSession())
	// Create an uploader with the session and default options
	downloader := s3manager.NewDownloader(sess)

	f, err := os.Create(c.String("file-name"))
	if err != nil {
		return cli.NewExitError("Failed to create file "+c.String("file-name"), 90)
	}
	if len(c.String("bucket")) == 0 {
		return cli.NewExitError("No bucket specified.", 87)
	}
	if len(c.String("prefix")) == 0 {
		return cli.NewExitError("No prefix specified.", 87)
	}
	// Write the contents of S3 Object to the file
	n, err := downloader.Download(f, &s3.GetObjectInput{
		Bucket: aws.String(c.String("bucket")),
		Key:    aws.String(c.String("prefix") + "/" + c.String("file-name")),
	})
	if err != nil {
		log.Fatal(err)
		return cli.NewExitError("Failed to download the file "+c.String("prefix")+" in the "+c.String("bucket")+" bucket.\n", 87)
	}
	fmt.Printf("file downloaded, %d bytes\n", n)
	return nil
}

package image

import (
	"fmt"
	"image"
	"image/png"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type Client struct {
	logger *log.Logger
}

func NewClient(logger *log.Logger) *Client {
	return &Client{
		logger: logger,
	}
}

func (c *Client) getClient() (*s3.S3, error) {
	config, err := session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(os.Getenv("SPACES_KEY"), os.Getenv("SPACES_SECRET"), ""),
		Endpoint:    aws.String(os.Getenv("SPACES_BUCKET")),
		Region:      aws.String("us-east-1"),
	})
	if err != nil {
		return nil, err
	}
	return s3.New(config), nil
}

func (c *Client) GetImage(imageId string) (*image.Image, error) {
	if os.Getenv("TEST") == "true" {
		var i image.Image
		return &i, nil
	}
	client, err := c.getClient()
	if err != nil {
		return nil, err
	}
	obj, err := client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(os.Getenv("SPACES_BUCKET")),
		Key:    aws.String(imageId),
	})
	if err != nil {
		return nil, err
	}
	i, err := png.Decode(obj.Body)
	if err != nil {
		return nil, err
	}
	return &i, nil
}

func (c *Client) PutImage(imageId string, i *image.Image) error {
	if os.Getenv("TEST") == "true" {
		return nil
	}
	c.logger.Println("getting client")
	client, err := c.getClient()
	if err != nil {
		return err
	}
	c.logger.Println("got client")
	c.logger.Println("list buckets")
	spaces, err := client.ListBuckets(nil)
	if err != nil {
		return err
	}
	for _, b := range spaces.Buckets {
		c.logger.Println(aws.StringValue(b.Name))
	}
	c.logger.Println("listed buckets")
	c.logger.Println("putting image")
	_, err = client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(os.Getenv("SPACES_BUCKET")),
		Key:    aws.String(imageId),
		Body:   strings.NewReader(imageId),
		ACL:    aws.String("private"),
		Metadata: map[string]*string{
			"x-amz-meta-my-key": aws.String("your-value"),
		},
	})
	if err != nil {
		fmt.Println(err.Error())
	}
	if err != nil {
		return err
	}
	c.logger.Println("put image")
	return nil
}

func (c *Client) DeleteImage(imageId string) error {
	if os.Getenv("TEST") == "true" {
		return nil
	}
	client, err := c.getClient()
	if err != nil {
		return err
	}
	_, err = client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(os.Getenv("SPACES_BUCKET")),
		Key:    aws.String(imageId),
	})
	if err != nil {
		return err
	}
	return nil
}

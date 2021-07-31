package image

import (
	"bytes"
	"image"
	"image/png"
	"log"
	"os"

	"github.com/minio/minio-go"
)

type Client struct {
	logger *log.Logger
}

func NewClient(logger *log.Logger) *Client {
	return &Client{
		logger: logger,
	}
}

func (c *Client) GetImage(imageId string) (*image.Image, error) {
	if os.Getenv("TEST") == "true" {
		var i image.Image
		return &i, nil
	}
	client, err := minio.New(os.Getenv("SPACES_URL"), os.Getenv("SPACES_API_KEY"), os.Getenv("SPACES_SECRET"), true)
	if err != nil {
		return nil, err
	}
	var obj *minio.Object
	obj, err = client.GetObject(os.Getenv("SPACES_BUCKET_NAME"), imageId, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	i, err := png.Decode(obj)
	if err != nil {
		return nil, err
	}
	return &i, nil
}

func (c *Client) PutImage(imageId string, i *image.Image) error {
	if os.Getenv("TEST") == "true" {
		return nil
	}
	client, err := minio.New(os.Getenv("SPACES_URL"), os.Getenv("SPACES_API_KEY"), os.Getenv("SPACES_SECRET"), true)
	if err != nil {
		return err
	}
	buffer := new(bytes.Buffer)
	err = png.Encode(buffer, *i)
	if err != nil {
		return err
	}
	_, err = client.PutObject(os.Getenv("SPACES_BUCKET_NAME"), imageId, buffer, -1, minio.PutObjectOptions{ContentType: "image/png"})
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) DeleteImage(imageId string) error {
	if os.Getenv("TEST") == "true" {
		return nil
	}
	client, err := minio.New(os.Getenv("SPACES_URL"), os.Getenv("SPACES_API_KEY"), os.Getenv("SPACES_SECRET"), true)
	if err != nil {
		return err
	}
	err = client.RemoveObject(os.Getenv("SPACES_BUCKET_NAME"), imageId)
	if err != nil {
		return err
	}
	return nil
}

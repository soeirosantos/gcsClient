package gcs

import (
	"context"
	"io"
	"log"
	"os"

	"google.golang.org/api/iterator"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

type Client struct {
	project string
	*storage.Client
}

var client *Client
var ctx = context.Background()

func GetClient(credentials, projectId string) *Client {
	if client == nil {
		c, err := storage.NewClient(ctx, option.WithCredentialsFile(credentials))
		if err != nil {
			log.Fatal(err)
		}
		client = &Client{projectId, c}
	}
	return client
}

func (c *Client) CreateBucket(bucketName string) error {
	bkt := c.Bucket(bucketName)
	return bkt.Create(ctx, c.project, nil)
}

func (c *Client) Upload(bucketName, filePath string) error {
	bkt := c.Bucket(bucketName)
	data, fileName := readFile(filePath)
	obj := bkt.Object(fileName)
	w := obj.NewWriter(ctx)
	_, err := w.Write(data)
	if err != nil {
		return err
	}
	if err := w.Close(); err != nil {
		return err
	}
	return nil
}

func (c *Client) Download(bucketName, file, location string) error {
	bkt := c.Bucket(bucketName)
	obj := bkt.Object(file)
	r, err := obj.NewReader(ctx)
	if err != nil {
		return err
	}
	defer r.Close()
	f, err := os.Create(location + string(os.PathSeparator) + file)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, r)
	return err
}

func (c *Client) Move(bucketName, from, to string) error {
	src := client.Bucket(bucketName).Object(from)
	dst := client.Bucket(bucketName).Object(to)
	if _, err := dst.CopierFrom(src).Run(ctx); err != nil {
		return err
	}
	if err := src.Delete(ctx); err != nil {
		return err
	}
	return nil
}

func (c *Client) List(bucketName string) ([]string, error) {
	bkt := c.Bucket(bucketName)
	it := bkt.Objects(ctx, nil)
	var files []string
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		files = append(files, attrs.Name)
	}
	return files, nil
}

func readFile(filePath string) ([]byte, string) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	fileinfo, err := file.Stat()
	if err != nil {
		log.Fatal(err)
	}
	filesize := fileinfo.Size()
	buffer := make([]byte, filesize)
	_, err = file.Read(buffer)
	if err != nil {
		log.Fatal(err)
	}
	return buffer, fileinfo.Name()
}

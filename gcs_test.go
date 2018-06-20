package gcs

import (
	"testing"
)

// It's an integration test. Replace the values below with your real values
const (
	ProjectId    = "project-id"
	BucketName   = "bucket-name"
	Credentials  = "/path/to/service_account.json"
	FileLocation = "./data/file_0_1"
)

func TestGetClient(t *testing.T) {
	client := GetClient(Credentials, ProjectId)
	if client == nil {
		t.Error("client should be not nil")
	}
}

func TestCreateBucket(t *testing.T) {
	client := GetClient(Credentials, ProjectId)
	err := client.CreateBucket(BucketName)
	if err != nil {
		t.Fatal(err)
	}
}

func TestReadFile(t *testing.T) {
	_, name := readFile(FileLocation)
	if name != "file_0_1" {
		t.Error("File was not properly read")
	}
}

func TestUpload(t *testing.T) {
	client := GetClient(Credentials, ProjectId)
	err := client.Upload(BucketName, FileLocation)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDownloald(t *testing.T) {
	client := GetClient(Credentials, ProjectId)
	err := client.Download(BucketName, "file_0_1", "./data/download")
	if err != nil {
		t.Fatal(err)
	}
}

func TestMove(t *testing.T) {
	client := GetClient(Credentials, ProjectId)
	const from = "file_0_1"
	const to = "moved/file_0_1"

	err := client.Move(BucketName, from, to)
	if err != nil {
		t.Fatal(err)
	}
}

func TestList(t *testing.T) {
	client := GetClient(Credentials, ProjectId)
	files, err := client.List(BucketName)
	if err != nil {
		t.Fatal(err)
	}
	if len(files) == 0 {
		t.Error("Files should have at least one element")
	}
}

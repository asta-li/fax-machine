package main

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"time"

	"cloud.google.com/go/storage"
	"golang.org/x/oauth2/google"
)

// TODO: set an lifecycle policy for pdfs and metadata that are left behind
// Uploads a file to GCS.
// See https://cloud.google.com/storage/docs/uploading-objects#storage-upload-object-go.
func uploadGCS(dataToWrite io.Reader, fileName string) error {
	// Create GCS connection
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	// Upload an object with storage.Writer.
	bucketName := os.Getenv("BUCKET_NAME")
	w := client.Bucket(bucketName).Object(fileName).NewWriter(ctx)
	if _, err = io.Copy(w, dataToWrite); err != nil {
		return fmt.Errorf("io.Copy: %v", err)
	}
	if err := w.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %v", err)
	}
	return nil
}

// Deletes an object from GCS.
func deleteGCS(fileName string) error {
	// Create GCS connection
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	// Upload an object with storage.Writer.
	bucketName := os.Getenv("BUCKET_NAME")
	obj := client.Bucket(bucketName).Object(fileName)
	if err := obj.Delete(ctx); err != nil {
		return fmt.Errorf("Object(%q).Delete: %v", fileName, err)
	}
	return nil
}

// Generates a signed URL that can be used to GET the given object.
func getSignedUrl(fileName string) (string, error) {
	bucketName := os.Getenv("BUCKET_NAME")
	credentials := os.Getenv("GCS_CREDENTIALS")
	jsonKey, err := ioutil.ReadFile(credentials)
	if err != nil {
		return "", fmt.Errorf("ioutil.ReadFile: %v", err)
	}
	conf, err := google.JWTConfigFromJSON(jsonKey)
	if err != nil {
		return "", fmt.Errorf("google.JWTConfigFromJSON: %v", err)
	}
	opts := &storage.SignedURLOptions{
		Scheme:         storage.SigningSchemeV4,
		Method:         "GET",
		GoogleAccessID: conf.Email,
		PrivateKey:     conf.PrivateKey,
		Expires:        time.Now().Add(10 * time.Minute),
	}
	signedUrl, err := storage.SignedURL(bucketName, fileName, opts)
	if err != nil {
		return "", fmt.Errorf("storage.SignedURL: %v", err)
	}
	return signedUrl, nil
}

func downloadGCS(filePath string) ([]byte, error) {

	bucketName := os.Getenv("BUCKET_NAME")
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	rc, err := client.Bucket(bucketName).Object(filePath).NewReader(ctx)
	if err != nil {
		return nil, fmt.Errorf("Object(%q).NewReader: %v", filePath, err)
	}
	defer rc.Close()

	data, err := ioutil.ReadAll(rc)
	if err != nil {
		return nil, fmt.Errorf("ioutil.ReadAll: %v", err)
	}
	log.Printf("Blob %v downloaded.\n", filePath)
	return data, nil
}

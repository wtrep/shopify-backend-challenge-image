package image

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"golang.org/x/oauth2/google"
	"io"
	"io/ioutil"
	"log"
	"os"
	"time"
)

const (
	urlExpirationDelay = 15 * time.Minute
)

// Upload the file to the specified GCP bucket at the object path
func uploadToBucket(data io.Reader, bucket, object string) error {
	client, err := storage.NewClient(context.Background())
	if err != nil {
		log.Println(err.Error())
		return err
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	wc := client.Bucket(bucket).Object(object).NewWriter(ctx)
	if _, err = io.Copy(wc, data); err != nil {
		log.Println(err.Error())
		return err
	}
	if err := wc.Close(); err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

// Generate a limited time download link for a specific object
func generateSignedURL(bucket, object string) (string, error) {
	serviceAccount := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	jsonKey, err := ioutil.ReadFile(serviceAccount)
	if err != nil {
		return "", fmt.Errorf("ioutil.ReadFile: %v", err)
	}

	conf, err := google.JWTConfigFromJSON(jsonKey)
	if err != nil {
		return "", fmt.Errorf("google.JWTConfigFromJSON: %v", err)
	}

	opts := &storage.SignedURLOptions{
		Method:         "GET",
		GoogleAccessID: conf.Email,
		PrivateKey:     conf.PrivateKey,
		Expires:        time.Now().Add(urlExpirationDelay),
	}
	u, err := storage.SignedURL(bucket, object, opts)
	if err != nil {
		return "", fmt.Errorf("storage.SignedURL: %v", err)
	}

	return u, nil
}

// Delete a specific file from a determined GCP bucket
func deleteFile(bucket, object string) error {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*15)
	defer cancel()

	o := client.Bucket(bucket).Object(object)
	if err := o.Delete(ctx); err != nil {
		return err
	}
	return nil
}

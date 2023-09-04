package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
)

func getBase64EncodedImageFromGCS(bucketName, objectName string) (string, error) {
	// Create a new GCS client.
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return "", fmt.Errorf("Failed to create GCS client: %v", err)
	}
	defer client.Close()

	// Get the bucket.
	bucket := client.Bucket(bucketName)

	// Get the object with the given name.
	obj := bucket.Object(objectName)

	// Read the contents of the object.
	reader, err := obj.NewReader(ctx)
	if err != nil {
		return "", err
	}
	defer reader.Close()
	data, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}

	fmt.Println("size of data: ", len(data), "B")

	if len(data) > 20971519 {
		fmt.Println("The image is too big. resizing...")
		data = resizeImage(data)
	}

	// Encode the object bytes to base64.
	encodedImg := base64.StdEncoding.EncodeToString(data)

	return encodedImg, nil
}

func listGCSFiles(bucketName string) ([]string, error) {
	// Create a new GCS client.
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("Failed to create GCS client: %v", err)
	}
	defer client.Close()

	// Get the bucket.
	bucket := client.Bucket(bucketName)

	// Get the bucket's objects with the given prefix.
	var files []string
	it := bucket.Objects(ctx, nil)
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("Failed to iterate over bucket objects: %v", err)
		}
		files = append(files, attrs.Name)
	}

	return files, nil
}

func uploadBytesToGCS(bucketName, dstFilename string, data []byte) error {
	// Create a new GCS client.
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("Failed to create GCS client: %v", err)
	}
	defer client.Close()

	// Create a new GCS object.
	bucket := client.Bucket(bucketName)
	dstObject := bucket.Object(dstFilename)

	// Write the source file to the GCS object.
	wc := dstObject.NewWriter(ctx)
	if _, err := wc.Write(data); err != nil {
		return fmt.Errorf("Failed to write file to GCS: %v", err)
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("Failed to close GCS writer: %v", err)
	}

	return nil
}

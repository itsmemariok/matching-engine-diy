package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"os"

	"github.com/disintegration/imaging"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func encodeImageToBase64(filePath string) (string, error) {
	// Read the image file.
	imgBytes, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	// Encode the image bytes to base64.
	encodedImg := base64.StdEncoding.EncodeToString(imgBytes)

	return encodedImg, nil
}

func createJSONLFile(bucketName, dstFilename string, items []Item) error {
	var buffer bytes.Buffer
	// Write each item to the file as a JSONL record.
	writer := bufio.NewWriter(&buffer)
	for _, item := range items {
		jsonBytes, err := json.Marshal(item)
		if err != nil {
			return fmt.Errorf("Failed to marshal JSON: %v", err)
		}
		jsonlRecord := fmt.Sprintf("%s\n", string(jsonBytes))

		_, err = writer.WriteString(jsonlRecord)
		if err != nil {
			return fmt.Errorf("Failed to write JSONL record: %v", err)
		}
	}
	writer.Flush()
	uploadBytesToGCS(bucketName, dstFilename, buffer.Bytes())

	return nil
}

func getBearerToken() string {
	ctx := context.Background()

	// Use Application Default Credentials (ADC) to authenticate the client.
	client, err := google.DefaultClient(ctx, "https://www.googleapis.com/auth/cloud-platform")
	if err != nil {
		panic(err)
	}

	// Get the access token from the client's token source.
	tokenSource := client.Transport.(*oauth2.Transport).Source
	token, err := tokenSource.Token()
	if err != nil {
		panic(err)
	}

	// Print the access token.
	return token.AccessToken
}

func resizeImage(data []byte) []byte {
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		panic(err)
	}

	maxSize := 20971519
	resizedImg := img
	for len(data) > maxSize {
		resizedImg = imaging.Resize(resizedImg, resizedImg.Bounds().Dx()/2, resizedImg.Bounds().Dy()/2, imaging.Lanczos)
		buf := new(bytes.Buffer)
		err := imaging.Encode(buf, resizedImg, imaging.JPEG)
		if err != nil {
			panic(err)
		}
		data = buf.Bytes()
	}
	return data
}

func deleteLocalFile(path string) {
	if fileExists(path) {
		err := os.Remove(path)
		if err != nil {
			panic(err)
		}
	}
}

func fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}

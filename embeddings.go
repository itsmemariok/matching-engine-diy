package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func getImageVector(project string, base64EncodedImg string) ([]float64, error) {
	// Create the JSON request body.
	bytesBase64Encoded := Image{
		BytesBase64Encoded: base64EncodedImg,
	}
	instance := Instance{
		Image: bytesBase64Encoded,
	}
	body := InstancesRequest{
		Instances: []Instance{
			instance,
		},
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		panic(err)
	}

	// Create the HTTP request.
	url := fmt.Sprintf("https://us-central1-aiplatform.googleapis.com/v1/projects/%s/locations/us-central1/publishers/google/models/multimodalembedding@001:predict", project)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		panic(err)
	}
	token := getBearerToken()
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	// Send the HTTP request and print the response.
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var res map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&res)
		if err != nil {
			panic(err)
		}
		fmt.Println(res)
		err = fmt.Errorf("Failed to get image vector: %v", err)
		return nil, err
	}

	var res ImagePredictionResponse
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		panic(err)
	}

	fmt.Println(resp.Status)
	// Print the predictions and deployed model ID.
	return res.Predictions[0].ImageEmbedding, nil
}

func getTextVector(project string, text string) ([]float64, error) {
	// Create the JSON request body.

	instance := InstanceText{
		Text: text,
	}
	body := TextInstancesRequest{
		Instances: []struct {
			Text string `json:"text"`
		}{
			{
				Text: instance.Text,
			},
		},
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		panic(err)
	}

	// Create the HTTP request.
	url := fmt.Sprintf("https://us-central1-aiplatform.googleapis.com/v1/projects/%s/locations/us-central1/publishers/google/models/multimodalembedding@001:predict", project)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		panic(err)
	}
	token := getBearerToken()
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	// Send the HTTP request and print the response.
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var res map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&res)
		if err != nil {
			panic(err)
		}
		fmt.Println(res)
		err = fmt.Errorf("Failed to get image vector: %v", err)
		return nil, err
	}

	var res TextPredictionResponse
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		panic(err)
	}

	fmt.Println(resp.Status)
	// Print the predictions and deployed model ID.
	return res.Predictions[0].TextEmbedding, nil
}

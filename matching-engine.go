package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func createMatchingEngineIndex(gcsbucket, project, location string) (string, error) {
	// Create the JSON request body.

	treeAhConfigParameters := TreeAhConfigParameters{
		LeafNodeEmbeddingCount:   5000,
		LeafNodesToSearchPercent: 3,
	}

	treeAhConfig := TreeAhConfig{
		TreeAhConfig: treeAhConfigParameters,
	}

	config := Config{
		Dimensions:                1408,
		ApproximateNeighborsCount: 100,
		ShardSize:                 "SHARD_SIZE_SMALL",
		DistanceMeasureType:       "DOT_PRODUCT_DISTANCE",
		AlgorithmConfig:           treeAhConfig,
	}

	metadata := Metadata{
		ContentsDeltaUri: "gs://" + gcsbucket + "/jsonl/",
		Config:           config,
	}

	searchIndex := SearchIndex{
		DisplayName:       "matching-engine-index",
		Metadata:          metadata,
		IndexUpdateMethod: "STREAM_UPDATE",
	}

	jsonBody, err := json.Marshal(searchIndex)
	if err != nil {
		panic(err)
	}

	// Create the HTTP request.
	url := "https://" + location + "-aiplatform.googleapis.com/v1/projects/" + project + "/locations/" + location + "/indexes"
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
		return "", err
	}
	var res IndexRelatedTaskResult
	json.NewDecoder(resp.Body).Decode(&res)

	parts := strings.Split(res.Name, "/operations/")

	return parts[0], nil
}

func createMatchingEngineIndexEndpoint(location, project, name string, isPublic bool) (string, error) {
	// Create the JSON request body.
	endpoint := Endpoint{
		DisplayName:           name,
		PublicEndpointEnabled: isPublic,
	}

	jsonBody, err := json.Marshal(endpoint)
	if err != nil {
		panic(err)
	}

	// Create the HTTP request.
	url := "https://" + location + "-aiplatform.googleapis.com/v1/projects/" + project + "/locations/" + location + "/indexEndpoints"
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
		return "", err
	}
	var res IndexRelatedTaskResult
	json.NewDecoder(resp.Body).Decode(&res)

	parts := strings.Split(res.Name, "/operations/")

	return parts[0], nil
}

func deployIndexToEndpoint(project, location, name, indexName, endpointName, machineType string) error {
	machineSpec := MachineSpec{
		MachineType: machineType,
	}

	dedicatedResource := DedicatedResource{
		MachineSpec:     machineSpec,
		MinReplicaCount: 1,
	}

	deployedIndex := DeployedIndex{
		ID:                 name,
		Index:              indexName,
		DedicatedResources: dedicatedResource,
	}

	deployIndex := DeployIndex{
		DeployedIndex: deployedIndex,
	}

	jsonBody, err := json.Marshal(deployIndex)
	if err != nil {
		panic(err)
	}

	// Create the HTTP request.
	url := "https://" + location + "-aiplatform.googleapis.com/v1/projects/" + project + "/locations/" + location + "/indexEndpoints/" + endpointName + ":deployIndex"
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
		fmt.Println(resp.Body)
		fmt.Println(resp.StatusCode)
		var res map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&res)
		if err != nil {
			panic(err)
		}
		fmt.Println(res)
		return err
	}
	var res map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&res)

	fmt.Println(res)

	return nil

}

func getDeployedIndexEndpointURL(project, location, endpointName string) string {
	// Create the HTTP request.
	url := "https://" + location + "-aiplatform.googleapis.com/v1/projects/" + project + "/locations/" + location + "/indexEndpoints/" + endpointName
	req, err := http.NewRequest(http.MethodGet, url, nil)
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

	var res IndexEndpoint
	json.NewDecoder(resp.Body).Decode(&res)

	return res.PublicEndpointDomainName
}

func search(endpointURL, indexEndpointID, deployedIndexID, projectNumber, location string, textVector []float64) error {

	datapoint := Datapoint{
		DatapointID:   "0",
		FeatureVector: textVector,
	}

	query := Query{
		Datapoint:     datapoint,
		NeighborCount: 5,
	}

	queryRequest := QueryRequest{
		DeployedIndexID: deployedIndexID,
		Queries:         []Query{query},
	}

	jsonBody, err := json.Marshal(queryRequest)
	if err != nil {
		panic(err)
	}

	// Create the HTTP request.
	url := "https://" + endpointURL + "/v1/projects/" + projectNumber + "/locations/" + location + "/indexEndpoints/" + indexEndpointID + ":findNeighbors"
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
		fmt.Println(resp.Body)
		fmt.Println(resp.StatusCode)
		var res map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&res)
		if err != nil {
			panic(err)
		}
		fmt.Println(res)
		return err
	}

	var res NearestNeighborsResponse
	json.NewDecoder(resp.Body).Decode(&res)

	cfg := readConfig()
	bucket := cfg.Bucket

	for _, neighbor := range res.NearestNeighbors[0].Neighbors {
		fmt.Println("filename: https://storage.cloud.google.com/" + bucket + "/" + neighbor.Datapoint.DatapointID)
		fmt.Println("distance: " + strconv.FormatFloat(neighbor.Distance, 'f', -1, 64))
		fmt.Println("")
	}

	return nil

}

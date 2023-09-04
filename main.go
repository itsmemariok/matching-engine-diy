package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	cfg := readConfig()
	var str string
	fmt.Println("Please choose the mode:")
	fmt.Println("1. Create Vector data, create Index and create Index Endpoint")
	fmt.Println("2. Deploy Index to Index Endpoint")
	fmt.Println("3. Get deployed Index Endpoint endpoint URL")
	fmt.Println("4. Search images using text")
	fmt.Println("")
	fmt.Print("Your choice?: ")
	fmt.Scan(&str)
	fmt.Println("")
	if str != "1" && str != "2" && str != "3" && str != "4" {
		fmt.Println("Invalid input")
		return
	}
	if str == "1" {
		createVectorData(cfg)
	}
	if str == "2" {
		deployIndex(cfg)
	}
	if str == "3" {
		getEndpointURL(cfg)
	}
	if str == "4" {
		getSearchResult(cfg)
	}
	return
}

func createVectorData(cfg GlobalConfig) {
	project := cfg.Project
	bucket := cfg.Bucket
	location := cfg.Location
	indexEndpointName := cfg.IndexEndpointName

	fmt.Println("Getting GCS files list...")
	objNames, err := listGCSFiles(bucket)

	fmt.Println("Inserting list to sqlite...")
	if err != nil {
		panic(err)
	}
	if err := insertFileNamesIntoDB(objNames); err != nil {
		panic(err)
	}

	fmt.Println("Generating Embedding...")
	var items []Item
	for _, objName := range objNames {
		if !strings.Contains(objName, ".jpg") {
			fmt.Println("Skipping: " + objName)
			continue
		}
		base64EncodedImg, _ := getBase64EncodedImageFromGCS(bucket, objName)
		fmt.Println("Generating embedding for: " + objName)
		imageVector, err := getImageVector(project, base64EncodedImg)
		if err != nil {
			fmt.Println("Failed to generate embedding for: " + objName + " So skipping it.")
			continue
		}
		items = append(items, Item{
			ID:        objName,
			Embedding: imageVector,
		})

	}

	fmt.Println("Creating JSONL file...")
	if err := createJSONLFile(bucket, "jsonl/items.json", items); err != nil {
		panic(err)
	}

	fmt.Println("creating Index...")
	indexName, _ := createMatchingEngineIndex(bucket, project, location)

	fmt.Println("creating Index Entrypoint...")
	createdIndexEndpointName, _ := createMatchingEngineIndexEndpoint(location, project, indexEndpointName, true)

	writeIndexNamesConfig(indexName, createdIndexEndpointName)

	fmt.Println("Done! Visit https://console.cloud.google.com/vertex-ai/matching-engine/indexes?project=" + project + " to see the progress of creating index. It should take about 1.5 hours to create index.")
	fmt.Println("After creating the index, run this program again and proceed to 2. Deploy Index to Index Endpoint.")
}

func deployIndex(cfg GlobalConfig) {
	project := cfg.Project
	location := cfg.Location
	machineType := cfg.IndexMachineType
	deployedIndexName := cfg.DeployedIndexName

	indexName, indexEndpointID := readIndexNamesConfig()
	fmt.Println("Deploying Index...")
	deployIndexToEndpoint(project, location, deployedIndexName, indexName, indexEndpointID, machineType)

	fmt.Println("Done! Visit https://console.cloud.google.com/vertex-ai/locations/us-central1/index-endpoints/" + indexEndpointID + "?project=" + project + " to see the progress of deploying index to index endpoint. It should take about 40 min to 1 hour to deploy.")
	fmt.Println("After deploying the index to index endpoint, run this program again and proceed to 3. Get deployed Index Endpoint endpoint URL.")
}

func getEndpointURL(cfg GlobalConfig) {
	project := cfg.Project
	location := cfg.Location

	_, indexEndpointID := readIndexNamesConfig()
	endpointURL := getDeployedIndexEndpointURL(project, location, indexEndpointID)
	writeDeployedIndexEndpointURL(endpointURL)

	fmt.Println("Done! Please run this program again and proceed to 4. Search images using text.")
}

func getSearchResult(cfg GlobalConfig) {
	project := cfg.Project
	projectNumber := cfg.ProjectNumber
	location := cfg.Location
	deployedIndexname := cfg.DeployedIndexName

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Search word?: ")
	scanner.Scan()
	text := scanner.Text()
	fmt.Println("")
	fmt.Println("You wrote: " + text)

	endpointURL := readDeployedIndexEndpointURL()
	_, indexEndpointID := readIndexNamesConfig()

	textVector, _ := getTextVector(project, text)

	search(endpointURL, indexEndpointID, deployedIndexname, projectNumber, location, textVector)
}

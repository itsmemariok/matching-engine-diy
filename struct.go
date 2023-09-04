package main

type Image struct {
	BytesBase64Encoded string `json:"bytesBase64Encoded"`
}

type Instance struct {
	Image Image `json:"image"`
}

type InstancesRequest struct {
	Instances []Instance `json:"instances"`
}

type ImagePrediction struct {
	ImageEmbedding []float64 `json:"imageEmbedding"`
}

type ImagePredictionResponse struct {
	Predictions     []ImagePrediction `json:"predictions"`
	DeployedModelId string            `json:"deployedModelId"`
}

type Item struct {
	ID        string    `json:"id"`
	Embedding []float64 `json:"embedding"`
}

type SearchIndex struct {
	DisplayName       string   `json:"display_name"`
	Metadata          Metadata `json:"metadata"`
	IndexUpdateMethod string   `json:"indexUpdateMethod"`
}

type Metadata struct {
	ContentsDeltaUri string `json:"contentsDeltaUri"`
	Config           Config `json:"config"`
}

type Config struct {
	Dimensions                int          `json:"dimensions"`
	ApproximateNeighborsCount int          `json:"approximateNeighborsCount"`
	ShardSize                 string       `json:"shardSize"`
	DistanceMeasureType       string       `json:"distanceMeasureType"`
	AlgorithmConfig           TreeAhConfig `json:"algorithmConfig"`
}

type TreeAhConfig struct {
	TreeAhConfig TreeAhConfigParameters `json:"treeAhConfig"`
}

type TreeAhConfigParameters struct {
	LeafNodeEmbeddingCount   int `json:"leafNodeEmbeddingCount"`
	LeafNodesToSearchPercent int `json:"leafNodesToSearchPercent"`
}

type Endpoint struct {
	DisplayName           string `json:"display_name"`
	PublicEndpointEnabled bool   `json:"publicEndpointEnabled"`
}

type DeployIndex struct {
	DeployedIndex DeployedIndex `json:"deployedIndex"`
}

type DeployedIndex struct {
	ID                 string            `json:"id"`
	Index              string            `json:"index"`
	DedicatedResources DedicatedResource `json:"dedicatedResources"`
}

type DedicatedResource struct {
	MachineSpec     MachineSpec `json:"machineSpec"`
	MinReplicaCount int         `json:"minReplicaCount"`
}

type MachineSpec struct {
	MachineType string `json:"machineType"`
}

type IndexRelatedTaskResult struct {
	Metadata struct {
		Type            string `json:"@type"`
		GenericMetadata struct {
			CreateTime string `json:"createTime"`
			UpdateTime string `json:"updateTime"`
		} `json:"genericMetadata"`
	} `json:"metadata"`
	Name string `json:"name"`
}

type AutomaticResources struct {
	MinReplicaCount int `json:"minReplicaCount"`
	MaxReplicaCount int `json:"maxReplicaCount"`
}

type DeployedIndexes struct {
	ID                 string             `json:"id"`
	Index              string             `json:"index"`
	DisplayName        string             `json:"displayName"`
	CreateTime         string             `json:"createTime"`
	IndexSyncTime      string             `json:"indexSyncTime"`
	AutomaticResources AutomaticResources `json:"automaticResources"`
	DeploymentGroup    string             `json:"deploymentGroup"`
}

type IndexEndpoint struct {
	Name                     string            `json:"name"`
	DisplayName              string            `json:"displayName"`
	DeployedIndexes          []DeployedIndexes `json:"deployedIndexes"`
	Etag                     string            `json:"etag"`
	CreateTime               string            `json:"createTime"`
	UpdateTime               string            `json:"updateTime"`
	PublicEndpointDomainName string            `json:"publicEndpointDomainName"`
}

type InstanceText struct {
	Text string `json:"text"`
}

type TextInstancesRequest struct {
	Instances []struct {
		Text string `json:"text"`
	} `json:"instances"`
}

type TextPredictionResponse struct {
	Predictions []struct {
		TextEmbedding []float64 `json:"textEmbedding"`
	} `json:"predictions"`
	DeployedModelID string `json:"deployedModelId"`
}

type Datapoint struct {
	DatapointID   string    `json:"datapoint_id"`
	FeatureVector []float64 `json:"feature_vector"`
}

type Query struct {
	Datapoint     Datapoint `json:"datapoint"`
	NeighborCount int       `json:"neighbor_count"`
}

type QueryRequest struct {
	DeployedIndexID string  `json:"deployed_index_id"`
	Queries         []Query `json:"queries"`
}

type CrowdingTag struct {
	CrowdingAttribute int64 `json:"crowdingAttribute"`
}

type DatapointResult struct {
	CrowdingTag CrowdingTag `json:"crowdingTag"`
	DatapointID string      `json:"datapointId"`
}

type Neighbor struct {
	Datapoint DatapointResult `json:"datapoint"`
	Distance  float64         `json:"distance"`
}

type NearestNeighbor struct {
	ID        int64      `json:"id"`
	Neighbors []Neighbor `json:"neighbors"`
}

type NearestNeighborsResponse struct {
	NearestNeighbors []NearestNeighbor `json:"nearestNeighbors"`
}

type GlobalConfig struct {
	Project           string `toml:"project"`
	ProjectNumber     string `toml:"projectnumber"`
	Location          string `toml:"location"`
	Bucket            string `toml:"bucket"`
	IndexMachineType  string `toml:"indexmachinetype"`
	IndexEndpointName string `toml:"indexendpointname"`
	DeployedIndexName string `toml:"deployedindexname"`
}

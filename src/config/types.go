package config

import (
	"time"
)

type ClusterInfo struct {
	Links []struct {
		Href string `json:"href"`
		Rel  string `json:"rel"`
	} `json:"links"`
	Clusters   []Cluster `json:"results"`
	TotalCount int       `json:"totalCount"`
}

type Cluster struct {
	AutoScaling struct {
		AutoIndexingEnabled bool `json:"autoIndexingEnabled"`
		Compute             struct {
			Enabled          bool `json:"enabled"`
			ScaleDownEnabled bool `json:"scaleDownEnabled"`
		} `json:"compute"`
		DiskGBEnabled bool `json:"diskGBEnabled"`
	} `json:"autoScaling"`
	BackupEnabled bool `json:"backupEnabled"`
	BiConnector   struct {
		Enabled        bool   `json:"enabled"`
		ReadPreference string `json:"readPreference"`
	} `json:"biConnector"`
	ClusterType       string `json:"clusterType"`
	ConnectionStrings struct {
		Standard    string `json:"standard"`
		StandardSrv string `json:"standardSrv"`
	} `json:"connectionStrings"`
	CreateDate               string  `json:"createDate"`
	DiskSizeGB               float64 `json:"diskSizeGB"`
	EncryptionAtRestProvider string  `json:"encryptionAtRestProvider"`
	GroupID                  string  `json:"groupId"`
	ID                       string  `json:"id"`
	Labels                   []struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	} `json:"labels"`
	Links []struct {
		Href string `json:"href"`
		Rel  string `json:"rel"`
	} `json:"links"`
	MongoDBMajorVersion   string    `json:"mongoDBMajorVersion"`
	MongoDBVersion        string    `json:"mongoDBVersion"`
	MongoURI              string    `json:"mongoURI"`
	MongoURIUpdated       time.Time `json:"mongoURIUpdated"`
	MongoURIWithOptions   string    `json:"mongoURIWithOptions"`
	Name                  string    `json:"name"`
	NumShards             int       `json:"numShards"`
	Paused                bool      `json:"paused"`
	PitEnabled            bool      `json:"pitEnabled"`
	ProviderBackupEnabled bool      `json:"providerBackupEnabled"`
	ProviderSettings      struct {
		ProviderName string `json:"providerName"`
		AutoScaling  struct {
			Compute struct {
				MaxInstanceSize interface{} `json:"maxInstanceSize"`
				MinInstanceSize interface{} `json:"minInstanceSize"`
			} `json:"compute"`
		} `json:"autoScaling"`
		BackingProviderName string `json:"backingProviderName"`
		RegionName          string `json:"regionName"`
		InstanceSizeName    string `json:"instanceSizeName"`
	} `json:"providerSettings"`
	ReplicationFactor int `json:"replicationFactor"`
	ReplicationSpec   struct {
		UsEast1 struct {
			AnalyticsNodes int `json:"analyticsNodes"`
			ElectableNodes int `json:"electableNodes"`
			Priority       int `json:"priority"`
			ReadOnlyNodes  int `json:"readOnlyNodes"`
		} `json:"US_EAST_1"`
	} `json:"replicationSpec"`
	ReplicationSpecs []struct {
		ID            string `json:"id"`
		NumShards     int    `json:"numShards"`
		RegionsConfig struct {
			UsEast1 struct {
				AnalyticsNodes int `json:"analyticsNodes"`
				ElectableNodes int `json:"electableNodes"`
				Priority       int `json:"priority"`
				ReadOnlyNodes  int `json:"readOnlyNodes"`
			} `json:"US_EAST_1"`
		} `json:"regionsConfig"`
		ZoneName string `json:"zoneName"`
	} `json:"replicationSpecs"`
	RootCertType         string `json:"rootCertType"`
	SrvAddress           string `json:"srvAddress"`
	StateName            string `json:"stateName"`
	VersionReleaseSystem string `json:"versionReleaseSystem"`
}

type Aggregation struct {
	DataSource string `json:"dataSource" :"dataSource"`
	Database   string `json:"database" :"database"`
	Collection string `json:"collection" :"collection"`
	Pipeline   string `bson:"pipeline" :"pipeline"`
}

type DatabaseList struct {
	Databases []Database `bson:"databases,omitempty"`
}

type Database struct {
	Name        *string `bson:"name,omitempty"`
	SizeOnDisk  int64   `bson:"sizeOnDisk,omitempty"`
	Empty       bool    `bson:"empty,omitempty"`
	Collections []string
}

type GridFsAggregation struct {
	ContentType    string `bson:"_id,omitempty"`
	FileCount      int64  `bson:"fileCount,omitempty"`
	TotalSize      int64  `bson:"totalSize,omitempty"`
	CollectionName string
}

type Document struct {
	ContentType    string `json:"_id,omitempty"`
	FileCount      int64  `json:"fileCount,omitempty"`
	TotalSize      int64  `json:"totalSize,omitempty"`
	CollectionName string
}

type AggregationOutput struct {
	Documents []Document `json:"documents,omitempty"`
}

type ResultSet struct {
	Documents []map[string]interface{} `json:"documents,omitempty"`
}

//db.fs.files.aggregate([   {$group:{_id:{DOMAIN_ID:"$DOMAIN_ID", contentType:"$contentType"},totalSize:{$sum:"$length"}, fileCount:{$sum: 1}}} ])

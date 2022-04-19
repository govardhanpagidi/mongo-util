package config

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	digest "github.com/mongodb-forks/digest"

	"net/http"
	"os"
)

type Config struct {
	Mongo    Mongo `json:"mongo,omitempty"`
	GCP      GCP   `json:"gcp,omitempty"`
	Interval int64 `json:"interval,omitempty"`
}

type Mongo struct {
	ProjectID     string `json:"project_id,omitempty"`
	ProjectName   string `json:"project_name,omitempty"`
	PublicKey     string `json:"pub_key,omitempty"`
	PrivateKey    string `json:"private_key,omitempty"`
	AtlasEndPoint string `json:"atlas_end_point,omitempty"`
}

type GCP struct {
	ProjectID string `json:"project_id,omitempty"`
	// Prefix		string `json:"prefix,omitempty"`
	ConfigPath string `json:"config_path,omitempty"`
}

type UserData struct {
	Users []MongoUser `json:"results,omitempty"`
}

type MongoUser struct {
	Username    string `json:"username,omitempty"`
	DBName      string `json:"databaseName,omitempty"`
	ProjectID   string `json:"groupId"`
	ProjectName string `json:"groupName"`
}

type Project struct {
	ID    string `json:"id,omitempty"`
	OrgID string `json:"orgId,omitempty"`
}

//LoadConfig from json file
func LoadConfig(path string, config *Config) (*Config, error) {

	configFile, err := os.Open(path)
	defer configFile.Close()

	if err != nil {
		return config, err
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)
	return config, err
}

//HttpCall is a Generic HTTP client util method
func HttpCall(method, uri string, payload []byte, config Mongo) (*http.Response, error) {

	transport := digest.NewTransport(config.PublicKey, config.PrivateKey)
	req, err := http.NewRequest(method, uri, bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	res, err := transport.RoundTrip(req)
	return res, err
}

//RandomString generates Random string, base64 encoding
func RandomString(length int) string {
	var Rando = rand.Reader
	b := make([]byte, length)
	_, _ = Rando.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)
}

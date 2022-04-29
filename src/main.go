package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	configuration "mongo-util/config"
	"os"
)

var config configuration.Config

func init() {

	//This is for development purpose, config.json is expected in the same dir.
	if _, err := configuration.LoadConfig("../src/config.json", &config); err != nil {
		log.Println("WARNING: config.json not found in the same directory:", err)
		//Load configuration file, path can be changed according the file where it exists
		if _, err := configuration.LoadConfig("/etc/config.json", &config); err != nil {
			log.Fatalln("config.json config load error:", err)
		}
	}

	//set gpc config
	err := os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", config.GCP.ConfigPath)
	if err != nil {
		log.Fatalln("GCP config load error:", err)
	}
	log.Println("config loaded successfully...")

}

const (
	Command          = "command"
	GridFSReport     = "gridfs_report"
	UpdatePasswords  = "update_passwords"
	ClusterReport    = "cluster_report"
	Execute          = "execute_query"
	Help             = "help"
	ProjectName      = "project_name"
	Cluster          = "cluster"
	Database         = "db"
	Collection       = "collection"
	AtlasPubKey      = "atlas_pub_key"
	AtlasPrivateKey  = "atlas_private_key"
	DataApiKey       = "data_api_key"
	ConnectionString = "connection_string"
	Query            = "query"
)

func main() {

	//command line arguments
	command := flag.String(Command, ClusterReport, "Operation name") //Default command is genreate-clusterdetails
	projectName := flag.String(ProjectName, "", "project name")
	clusterName := flag.String(Cluster, "", "database name")
	dbName := flag.String(Database, "", "database name")
	collName := flag.String(Collection, "", "collection name")
	pubKey := flag.String(AtlasPubKey, "", "atlas public key")
	privateKey := flag.String(AtlasPrivateKey, "", "atlas private key")
	dataApiKey := flag.String(DataApiKey, "", "data api key")
	connString := flag.String(ConnectionString, "", "Mongodb connection string")
	query := flag.String(Query, "", "query to execute")

	flag.Parse()

	if command == nil {
		log.Fatalln("enter the operation name with -command argument. eg: -command=clusterreport ")
		return
	}

	log.Println("entered command:", *command)
	//Read the commands and execute respective operation based on arguments
	switch *command {
	case UpdatePasswords:
		if err := setAtlasConfig(pubKey, privateKey); err != nil {
			log.Println(err)
			return
		}
		if err := updatePasswords(projectName); err != nil {
			log.Println(err)
			return
		}
	case GridFSReport:
		//Generate the documents/files details as a CSV report
		if err := setAggregationConfig(clusterName, dbName, collName, dataApiKey, connString); err != nil {
			log.Println(err)
			return
		}
		if err := executeGridFSQuery(config.Mongo, clusterName, dbName, collName); err != nil {
			log.Println(err)
		}

		return
	case ClusterReport:
		if err := setAtlasConfig(pubKey, privateKey); err != nil {
			log.Println(err)
			return
		}
		if err := generateClusterDetailReport(projectName); err != nil {
			log.Println(err)
			return
		}

	case Execute:
		if err := setAggregationConfig(clusterName, dbName, collName, dataApiKey, connString); err != nil {
			log.Println(err)
			return
		}

		if err := executeQuery(query); err != nil {
			log.Println("main error:", err)
		}
	default:
		log.Fatalln("no such command ", *command)
	}
	return
}

func setAtlasConfig(pubKey, privateKey *string) error {

	if *pubKey == "" || *privateKey == "" {
		return errors.New(fmt.Sprintf("%s or %s is missing, please check", AtlasPubKey, AtlasPrivateKey))
	}
	config.Mongo.PublicKey = *pubKey
	config.Mongo.PrivateKey = *privateKey
	return nil
}

func setAggregationConfig(clusterName, dbName, collName, dataApiKey, connString *string) error {

	if connString != nil && *connString != "" {
		config.Mongo.ConnectionString = connString
	}
	if config.Mongo.ConnectionString == nil || *config.Mongo.ConnectionString == "" {
		return errors.New(fmt.Sprintf("connection string is requried"))
	}

	if dataApiKey != nil {
		config.Mongo.ApiKey = *dataApiKey
	}

	//Validation for query params
	if clusterName == nil || *clusterName == "" {
		return errors.New(fmt.Sprintf("invalid request: %s parameter required to generate GridFS report", Cluster))
	}
	return nil
}

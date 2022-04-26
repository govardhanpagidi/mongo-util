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
	Command         = "command"
	GridFSReport    = "gridfs_report"
	UpdatePasswords = "update_passwords"
	ClusterReport   = "cluster_report"
	Execute         = "execute_query"
	Help            = "help"
	ProjectName     = "project_name"
	Cluster         = "cluster"
	Database        = "db"
	Collection      = "collection"
	AtlasPubKey     = "atlas_pub_key"
	AtlasPrivateKey = "atlas_private_key"
	DataApiKey      = "data_api_key"
	Query           = "query"
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
		if err := setAggregationConfig(clusterName, dbName, collName, dataApiKey); err != nil {
			log.Println(err)
			return
		}
		if err := executeGridFSQuery(config.Mongo, clusterName, dbName, collName, nil); err != nil {
			log.Println(err)
		}
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
		if err := setAggregationConfig(clusterName, dbName, collName, dataApiKey); err != nil {
			log.Println(err)
			return
		}

		if err := executeGridFSQuery(config.Mongo, clusterName, dbName, collName, query); err != nil {
			log.Println("main error:", err)
		}
	case Help:
		printHelpSection()
	default:
		log.Fatalln("no such command", command)
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

func setAggregationConfig(clusterName, dbName, collName, dataApiKey *string) error {
	if dataApiKey == nil || config.Mongo.DataEndPoint == "" {
		return errors.New(fmt.Sprintf("%s/data_end_point is missing, please check arguments or config.json", DataApiKey))
	}
	config.Mongo.ApiKey = *dataApiKey

	//Validation for query params
	if *clusterName == "" || *dbName == "" || *collName == "" {
		return errors.New(fmt.Sprintf("invalid request: %s, %s, %s parameters required to generate GridFS report", Cluster, Database, Collection))
	}
	return nil
}

func printHelpSection() {
	log.Println("-command	command name, possible values are cluserreport,gridfsreport and updatepasswords ")
	log.Println("-project-name	project name, get this value from your atlas dashboard")
	log.Println("-db	database name")
	log.Println("-collection	collection name")
	return
}

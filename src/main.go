package main

import (
	"flag"
	"fmt"
	"log"
	configuration "mongo-util/config"
	"mongo-util/gcp"
	"mongo-util/mongo"
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
	GridFSReport    = "gridfsreport"
	UpdatePasswords = "updatepasswords"
	ClusterReport   = "clusterreport"
)

func main() {

	command := flag.String("command", ClusterReport, "Operation name")
	projectName := flag.String("projectname", "zebra", "project name")
	flag.Parse()
	if command == nil {
		log.Fatalln("enter the operation name with -command argument. eg: -command=clusterreport ")
		return
	}

	//config.Mongo.ProjectName = args[0]
	log.Println(*command)
	switch *command {
	case GridFSReport:

		generateGridFsReport(config.Mongo, nil, nil)
	case UpdatePasswords:
		if projectName == nil {
			log.Fatalln("projectName argument is missing the value")
		}
		config.Mongo.ProjectName = *projectName
		project, err := mongo.GetProjectByProjectName(config.Mongo)
		if err != nil {
			log.Fatalln(err)
		}
		config.Mongo.ProjectID = project.ID
		updateMongoUsers()
	case ClusterReport:
		if projectName == nil {
			log.Fatalln("projectName argument is missing the value")
		}
		config.Mongo.ProjectName = *projectName
		//Read project name from args and get the project id
		project, err := mongo.GetProjectByProjectName(config.Mongo)
		if err != nil {
			log.Fatalln(err)
		}
		config.Mongo.ProjectID = project.ID
		err = generateClustersReport(config.Mongo)
		if err != nil {
			log.Fatalln("error:", err)
		}

	}

	//log.Printf("Project Details %+v", project)

	//Setting projectId to config

	//update with new password and push the same to GCP
	//updateMongoUsers()

	//GetCluserReports
	//generateClustersReport(config.Mongo)
	//mongo.ReadAggregation(config.Mongo)

}

func getReports() {
	//Fetch the list of mongodb users
	users, err := mongo.GetUsersByProject(config.Mongo)
	if err != nil {
		return
	}

	entries := make([][]string, len(users)+1) //+1 is for including headers in the CSV
	entries[0] = []string{"Username", "ProjectName", "ProjectId"}
	for ind, value := range users {
		entries[ind+1] = []string{value.Username, config.Mongo.ProjectName, value.ProjectID}
	}
	err = mongo.GenerateCSV(entries, "users")
	if err != nil {
		log.Println("Error:", err)
		return
	}
	log.Println("Reports generated successfully...")
	return
}

//ClusterColumns for quick reference  as we need to follow the same order
var ClusterColumns = []string{"Name", "GroupId", "ClusterType", "DiskSizeGB", "NumShards", "ReplicationFactor", "CreatedDate", "BackupEnabled", "mongoDBMajorVersion", "mongoDBVersion"}

func generateClustersReport(config configuration.Mongo) error {
	clusterInfo, err := mongo.GetClusterInfo(config)
	if err != nil || (clusterInfo == nil || len(clusterInfo.Clusters) <= 0) {
		return err
	}

	var clusterEntries = make([][]string, len(clusterInfo.Clusters)+1)
	clusterEntries[0] = ClusterColumns
	for ind, cluster := range clusterInfo.Clusters {
		clusterEntries[ind+1] = []string{
			cluster.Name,
			cluster.GroupID,
			cluster.ClusterType,
			fmt.Sprintf("%f", cluster.DiskSizeGB),
			fmt.Sprintf("%d", cluster.NumShards),
			fmt.Sprintf("%d", cluster.ReplicationFactor),
			cluster.CreateDate,
			fmt.Sprintf("%t", cluster.BackupEnabled),
			cluster.MongoDBMajorVersion,
			cluster.MongoDBVersion,
		}
	}

	err = mongo.GenerateCSV(clusterEntries, fmt.Sprintf("%s_clusterinfo_%s", config.ProjectName, configuration.TimeNow()))
	if err != nil {
		log.Println("Error:", err)
		return err
	}
	log.Println("Reports generated successfully...")
	return err
}

var dbColumns = []string{"Name", "GroupId", "ClusterType", "DiskSizeGB", "NumShards", "ReplicationFactor", "CreatedDate", "BackupEnabled", "mongoDBMajorVersion", "mongoDBVersion"}

func generateDatabaseReport(config configuration.Mongo) error {
	dbs, err := mongo.GetDatabaseInfo(config)
	if err != nil || dbs == nil {
		return err
	}

	var clusterEntries = make([][]string, len(*dbs)+1)
	clusterEntries[0] = dbColumns
	for ind, db := range *dbs {
		clusterEntries[ind+1] = []string{
			db.Name,
			fmt.Sprintf("%d", db.SizeOnDisk),
		}
	}

	err = mongo.GenerateCSV(clusterEntries, fmt.Sprintf("%s_clusterinfo_%s", config.ProjectName, configuration.TimeNow()))
	if err != nil {
		log.Println("Error:", err)
		return err
	}
	log.Println("Reports generated successfully...")
	return err
}

var gridfsColumns = []string{"Database", "Collection", "ContentType", "FileCount", "TotalSize"}

func generateGridFsReport(config configuration.Mongo, dbName, collectionName *string) error {

	var dbNames []string
	if dbName != nil {
		dbNames = append(dbNames, *dbName)
	} else {
		dbs, err := mongo.GetDatabaseInfo(config)
		for _, db := range *dbs {
			if err != nil || dbs == nil {
				return err
			}
			dbNames = append(dbNames, db.Name)
		}
	}

	//list of collectionNames with dbname as key pair
	//collections := make(map[string]string)
	//If projectName provided get the report for all the databases
	var fsEntries [][]string
	fsEntries = append(fsEntries, gridfsColumns)

	for _, dbName := range dbNames {
		gridFsList, err := mongo.GetGridFsInfoByDB(config, dbName, collectionName)
		if err != nil || gridFsList == nil {
			return err
		}
		//Append entries

		for _, data := range *gridFsList {
			fsEntries = append(fsEntries, []string{
				dbName,
				data.CollectionName,
				data.ContentType,
				fmt.Sprintf("%d", data.FileCount),
				fmt.Sprintf("%d", data.TotalSize),
			})
			//fmt.Printf("data :%s %s %s %s %s ", dbName,data.CollectionName, data.ContentType,	fmt.Sprintf("%d", data.FileCount),fmt.Sprintf("%d", data.TotalSize))
		}
	}

	//Generate CSV
	err := mongo.GenerateCSV(fsEntries, fmt.Sprintf("%s_FSINFO_%s", config.ProjectName, configuration.TimeNow()))
	if err != nil {
		log.Println("Error:", err)
		return err
	}
	log.Println("Report generated successfully...")
	return err
}

func generateAggregationReport(config configuration.Mongo) error {
	clusterInfo, err := mongo.GetClusterInfo(config)
	if err != nil && len(clusterInfo.Clusters) <= 0 {
		return err
	}

	var clusterEntries = make([][]string, len(clusterInfo.Clusters)+1)
	clusterEntries[0] = configuration.ClusterColumns
	for ind, cluster := range clusterInfo.Clusters {
		clusterEntries[ind+1] = []string{
			cluster.Name,
			cluster.GroupID,
			cluster.ClusterType,
			fmt.Sprintf("%f", cluster.DiskSizeGB),
			fmt.Sprintf("%d", cluster.NumShards),
			fmt.Sprintf("%d", cluster.ReplicationFactor),
			cluster.CreateDate,
			fmt.Sprintf("%t", cluster.BackupEnabled),
			cluster.MongoDBMajorVersion,
			cluster.MongoDBVersion,
		}
	}

	err = mongo.GenerateCSV(clusterEntries, fmt.Sprintf("%s_clusterinfo_%s", config.ProjectName, configuration.TimeNow()))
	if err != nil {
		log.Println("Error:", err)
		return err
	}
	log.Println("Reports generated successfully...")
	return err
}

func updateMongoUsers() {
	//Fetch the list of mongodb users
	users, err := mongo.GetUsersByProject(config.Mongo)
	if err != nil {
		log.Fatalln("get users:", err)
		return
	}
	log.Printf("Total users under %s : %d", config.Mongo.ProjectID, len(users))

	for _, userInfo := range users {

		//Length 16 is expected
		pwd := configuration.RandomString(16)
		//log.Printf("Updating user %s with new password %s", result.Username, pwd)
		if err := mongo.UpdatePassword(pwd, userInfo, config.Mongo); err != nil {
			log.Printf("unable to change %s password for the DB: %s", userInfo.Username, userInfo.DBName)
			//log.Println(err)
			continue
		}

		//Save to GCP secret manager
		if err := gcp.SaveSecret(config, userInfo, pwd); err != nil {
			log.Printf("gcp error: while updating the user %s for the DB %s :", userInfo.Username, userInfo.DBName)
			log.Println(err)
		}
		log.Printf("updated password for	 %s		%s", userInfo.Username, pwd)
	}
}

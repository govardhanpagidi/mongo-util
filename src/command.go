package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	configuration "mongo-util/config"
	gcp "mongo-util/gcp"
	mongo "mongo-util/mongo"
)

func updatePasswords(projectName *string) error {
	if projectName == nil {
		return errors.New("project_name argument is missing the value")
	}
	config.Mongo.ProjectName = *projectName
	//Get projectId for a given project name through Atlas API
	project, err := mongo.GetProjectByProjectName(config.Mongo)
	if err != nil {
		log.Fatalln("GetProjectByProjectName error :", err)
		return err
	}
	config.Mongo.ProjectID = project.ID
	err = updateMongoUsers()
	if err != nil {
		return err
	}
	log.Println("passwords update successful")
	return err
}

func generateClusterDetailReport(projectName *string) error {
	if projectName == nil {
		return errors.New("project-name argument is missing the value")
	}
	config.Mongo.ProjectName = *projectName

	//Get projectId for a given project name through Atlas API
	project, err := mongo.GetProjectByProjectName(config.Mongo)
	if err != nil {
		return err
	}
	config.Mongo.ProjectID = project.ID

	//Generate the cluster details as a CSV report
	err = generateClustersReport(config.Mongo)
	if err != nil {
		return err
	}
	log.Println("cluster report generated successfully")
	return err
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
	return err
}

var dbColumns = []string{"Name", "GroupId", "ClusterType", "DiskSizeGB", "NumShards", "ReplicationFactor", "CreatedDate", "BackupEnabled", "mongoDBMajorVersion", "mongoDBVersion"}

var gridfsColumns = []string{"Database", "Collection", "ContentType", "FileCount", "TotalSize"}

func updateMongoUsers() error {
	//Fetch the list of mongodb users
	users, err := mongo.GetUsersByProject(config.Mongo)
	if err != nil {
		log.Fatalln("get users:", err)
		return err
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
			return err
		}
	}
	return err
}

func executeQuery(query *string) error {
	// Just keep it for testing..
	//*query = `{
	//		  	"dataSource": "zebra",
	//		  	"database": "myFirstDatabase",
	//			"collection":"uploads.files",
	//		  	"pipeline":[{
	//						"$group" : {
	//								"_id": "$contentType",
	//								"totalSize": {
	//									"$sum": "$length"
	//								},
	//								"fileCount":{"$sum" : 1}
	//						}
	//					}
	//				]
	//			}`
	var jsonMap map[string]interface{}
	err := json.Unmarshal([]byte(*query), &jsonMap)
	if err != nil {
		return err
	}
	qryBytes, err := json.Marshal(jsonMap)
	if err != nil {
		return err
	}
	result, err := mongo.ExecuteQueryAPI(config.Mongo, string(qryBytes))
	if err != nil {
		return err
	}

	//Unmarshal with map
	var aggregationOutput configuration.ResultSet
	err = json.Unmarshal(result, &aggregationOutput)
	fmt.Println("output:", string(result))

	//Iterate the result-set
	var entries [][]string

	var keys []string

	entries = append(entries, keys)
	for ind, val := range aggregationOutput.Documents {
		if ind == 0 {
			for k := range aggregationOutput.Documents[ind] {
				keys = append(keys, k)
			}
			entries = append(entries, keys)
		}

		var values []string
		for key := range val {
			reflectValue := aggregationOutput.Documents[ind][key]
			values = append(values, fmt.Sprintf("%+v", reflectValue))

		}
		entries = append(entries, values)
	}

	//Generate CSV
	if err = mongo.GenerateCSV(entries, fmt.Sprintf("%s_Results_%s", config.Mongo.ProjectName, configuration.TimeNow())); err != nil {
		return err
	}
	return err
}

func executeGridFSQuery(config configuration.Mongo, clusterName, dbName, collectionName *string) error {

	databases, err := getDatabasesWithCollections(config, dbName, collectionName)

	var fsEntries [][]string
	fsEntries = append(fsEntries, gridfsColumns)
	for _, db := range databases {
		for _, collection := range db.Collections {

			//Make a driver call to run the aggregation
			gridFsMap, err := mongo.ExecuteQuery(config, dbName, &collection)
			if err != nil || gridFsMap == nil {
				log.Println("ExecuteQuery returned error/nil :", err)
				return err
			}

			row := []string{*db.Name, collection}
			//Append values
			for _, val := range *gridFsMap {
				row = append(row, fmt.Sprintf("%+v", val))
			}
			fsEntries = append(fsEntries, row)
		}

	}

	//Generate CSV
	err = mongo.GenerateCSV(fsEntries, fmt.Sprintf("%s_FSINFO_%s", config.ProjectName, configuration.TimeNow()))
	if err != nil {
		return err
	}
	log.Println("files report generated successfully")
	return err
}

func getDatabasesWithCollections(config configuration.Mongo, dbName, collectionName *string) (databases []configuration.Database, err error) {

	//If database name is not provided consider fetching all the databases and collections in a cluster
	if dbName == nil || *dbName == "" {
		dbs, err := mongo.GetDatabaseInfo(config)
		if err != nil || dbs == nil {
			//log.Println("GetDatabaseInfo error:", err)
			return nil, err
		}
		databases = *dbs
	} else { // if database is provided
		var collections []string
		//database provided but not collection
		if collectionName == nil || *collectionName == "" {
			colls, err := mongo.GetCollections(config, *dbName)
			if err != nil || colls == nil {
				return nil, err
			}
			for _, collection := range *colls {
				collections = append(collections, collection.(string))
			}
		} else { //database and collection provided
			collections = append(collections, *collectionName)
		}
		//append an entry into database array since dbname is provided
		databases = append(databases, configuration.Database{
			Name:        dbName,
			Collections: collections,
		})
	}
	return
}

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
		log.Println("error while generating the files reports:", err)
		return err
	}
	log.Println("files report generated successfully")
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

func executeQuery(query *string) {
	// Just keep it for testing..
	//*query = `{
	//		  	"dataSource": "zebra",
	//		  	"database": "myFirstDatabase",
	//			"collection":"uploads.files",
	//		  	"pipeline":[{
	//						"$group" : {
	//								"_id": {"_id":"$contentType"},
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
		log.Fatalln(err)
	}
	qryBytes, err := json.Marshal(jsonMap)
	if err != nil {
		log.Fatalln(err)
	}
	result, err := mongo.ExecuteQuery(config.Mongo, string(qryBytes))
	if err != nil {
		log.Println("error:", err)
		return
	}
	var aggregationOutput configuration.AggregationOutput
	err = json.Unmarshal(result, &aggregationOutput)
	fmt.Println("output:", string(result))

	//Generate CSV
	//var entries [][]string
	//count := 0
	//var keys []string
	//for k := range aggregationOutput.Documents {
	//	keys = append(keys, )
	//}
	//entries = append(entries, keys)
	//for _, val := range aggregationOutput.Documents {
	//	var values []string
	//	for i := 0; i < len(keys); i++ {
	//		keys =
	//	}
	//	entries = append(entries, values)
	//}
	return
}

func executeGridFSQuery(config configuration.Mongo, clusterName, dbName, collectionName, query *string) error {
	// Build query
	defaultQuery := `{
			  	"dataSource": "` + *clusterName + `",
			  	"database": "` + *dbName + `",
				"collection":"` + *collectionName + `",
			  	"pipeline":[{
							"$group" : {
									"_id": "$contentType",
									"totalSize": {
										"$sum": "$length"
									},
									"fileCount":{"$sum" : 1}
							}
						}
					]
				}`

	if query == nil || *query == "" {
		query = &defaultQuery
	}
	var jsonMap map[string]interface{}
	err := json.Unmarshal([]byte(*query), &jsonMap)
	if err != nil {
		log.Fatalln(err)
	}
	qryBytes, err := json.Marshal(jsonMap)
	if err != nil {
		log.Fatalln(err)
	}

	var dbNames []string
	if dbName != nil {
		dbNames = append(dbNames, *dbName)
	} else {
		dbs, err := mongo.GetDatabaseInfo(config)
		if err != nil {
			log.Println("GetDatabaseInfo error:", err)
			return err
		}
		for _, db := range *dbs {
			if err != nil || dbs == nil {
				log.Println("for loop GetDatabaseInfo error:", err)

				return err
			}
			dbNames = append(dbNames, db.Name)
		}
	}

	//If projectName provided get the report for all the databases
	var fsEntries [][]string
	fsEntries = append(fsEntries, gridfsColumns)

	for _, dbName := range dbNames {

		result, err := mongo.ExecuteQuery(config, string(qryBytes))
		if err != nil {
			log.Println("ExecuteQuery error:", err)
			return nil
		}
		fmt.Println("result:", string(result))
		var gridFsList configuration.AggregationOutput
		err = json.Unmarshal(result, &gridFsList)

		if err != nil || gridFsList.Documents == nil {
			return err
		}
		//Append entries

		for _, data := range gridFsList.Documents {
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
	if err = mongo.GenerateCSV(fsEntries, fmt.Sprintf("%s_FSINFO_%s", config.ProjectName, configuration.TimeNow())); err != nil {
		return err
	}
	log.Println("files report generated successfully")
	return err
}

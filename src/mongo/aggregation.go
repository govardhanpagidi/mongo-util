package mongo

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	configuration "mongo-util/config"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetGridFsInfoByDB(config configuration.Mongo, dbName string, collName *string) (*[]configuration.GridFsAggregation, error) {

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.ConnectionString))
	if err != nil {
		panic(err)
	}
	defer client.Disconnect(ctx)

	db := client.Database(dbName)
	var collNames []string

	if collName != nil {
		collNames = append(collNames, *collName)
	}
	if names, err := getCollections(db); err == nil && names != nil && len(*names) > 0 {
		//fmt.Printf("collections : %+v", names)
		for _, val := range *names {
			//fmt.Println(val.(string))
			collNames = append(collNames, val.(string))
		}
	}

	var gridFsData []configuration.GridFsAggregation
	for _, collName := range collNames {
		collection := db.Collection(collName)
		data, err := getGridFsAggregationData(ctx, collection, collName)
		if err == nil || data != nil {
			for _, d := range *data {
				d.CollectionName = collName
				gridFsData = append(gridFsData, d)
			}
		}
	}

	return &gridFsData, err
}

func getGridFsAggregationData(ctx context.Context, collection *mongo.Collection, collName string) (*[]configuration.GridFsAggregation, error) {
	groupStage :=
		bson.D{
			{"$group", bson.D{
				{"_id", "$contentType"},
				{"totalSize", bson.D{{"$sum", "$length"}}},
				{"fileCount", bson.D{{"$sum", 1}}}}}}

	cursor, err := collection.Aggregate(ctx, mongo.Pipeline{groupStage})
	var gridFsInfo []configuration.GridFsAggregation
	//var gridFsInfo []bson.M
	if cursor != nil {
		if err = cursor.All(ctx, &gridFsInfo); err != nil {
			fmt.Println("getGridFsAggregationData, cursor.All error:", err)
			return nil, err
		}
	}
	//fmt.Println(gridFsInfo)
	return &gridFsInfo, err
}

func getCollections(database *mongo.Database) (*[]interface{}, error) {

	result, err := database.ListCollectionNames(
		context.TODO(),
		bson.D{})

	if err != nil {
		log.Fatal(err)
	}

	collRes, err := json.MarshalIndent(result, "", "	")
	var collNames []interface{}
	if err := json.Unmarshal(collRes, &collNames); err != nil {
		fmt.Println("unmarshal err:", err)
		return nil, err
	}

	return &collNames, err
}

func getDatabaseInfo(config configuration.Mongo) (*[]configuration.Database, error) {

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.ConnectionString))
	if err != nil {
		panic(err)
	}
	defer client.Disconnect(ctx)

	db := client.Database("admin")
	command := bson.D{{"listDatabases", 1}}
	var res bson.M
	err = db.RunCommand(context.TODO(), command).Decode(&res)
	if err != nil {
		panic(err)
	}

	//fmt.Printf("%+v", &res)
	btRes, err := json.MarshalIndent(res, "", "	")

	var dbList configuration.DatabaseList
	//fmt.Println("b4 unmarshal:")
	if err := json.Unmarshal(btRes, &dbList); err != nil {
		fmt.Println("unmarshal err:", err)
		return nil, err
	}
	return &dbList.Databases, err

}

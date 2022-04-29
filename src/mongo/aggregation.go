package mongo

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	configuration "mongo-util/config"
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func getCollectionsByDB(db *mongo.Database) (*[]interface{}, error) {

	result, err := db.ListCollectionNames(
		context.TODO(),
		bson.D{})

	if err != nil {
		//log.Fatal(err)
		return nil, err
	}

	collRes, err := json.MarshalIndent(result, "", "	")
	var collNames []interface{}
	if err := json.Unmarshal(collRes, &collNames); err != nil {
		fmt.Println("unmarshal err:", err)
		return nil, err
	}

	return &collNames, err
}

func GetCollections(config configuration.Mongo, dbName string) (*[]interface{}, error) {

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(*config.ConnectionString))
	if err != nil {
		panic(err)
	}
	defer client.Disconnect(ctx)

	db := client.Database(dbName)
	return getCollectionsByDB(db)
}

func getDatabaseInfo(config configuration.Mongo) (*[]configuration.Database, error) {

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(*config.ConnectionString))
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

func runAggregateQuery(config configuration.Mongo, dbName, collName *string) (*map[string]interface{}, error) {

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(*config.ConnectionString))

	if err != nil {
		panic(err)
	}
	defer client.Disconnect(ctx)

	//{
	//	"$group" : {
	//	"_id": "$contentType",
	//		"totalSize": {
	//		"$sum": "$length"
	//	},
	//	"fileCount":{"$sum" : 1}
	//}
	//}
	aggregateStage := bson.D{
		bson.E{Key: "$group", Value: bson.D{
			bson.E{Key: "_id", Value: "$contentType"},
			bson.E{Key: "fileCount", Value: bson.D{primitive.E{Key: "$sum", Value: 1}}},
			bson.E{Key: "totalSize", Value: bson.D{primitive.E{Key: "$sum", Value: "$length"}}},
		}}}

	db := client.Database(*dbName)

	cursor, err := db.Collection(*collName).Aggregate(ctx, mongo.Pipeline{aggregateStage})
	if err != nil {
		fmt.Println("execute error :", err)
		return nil, err
	}

	//var gridFsres []interface{}
	result := make(map[string]interface{})
	for cursor.Next(ctx) {
		var gridFsData interface{}
		if err = cursor.Decode(&gridFsData); err != nil {
			log.Println("cursor error:", err)
		}

		entry := gridFsData.(primitive.D)
		//var intVal interface{}
		for _, v := range entry {
			log.Println(v.Key, v.Value)

			if isNil(v.Value) {
				// intVal = " "
				result[v.Key] = ""
				continue
			}
			result[v.Key] = v.Value
		}
	}
	return &result, err
}

func isNil(i interface{}) bool {
	//fmt.Println("kind:", reflect.ValueOf(i).Kind())
	return i == nil || (reflect.ValueOf(i).Kind() == reflect.Invalid) || (reflect.ValueOf(i).Kind() == reflect.Ptr && reflect.ValueOf(i).IsNil())
}

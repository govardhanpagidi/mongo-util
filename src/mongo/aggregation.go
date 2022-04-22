package mongo

import (
	"context"
	"encoding/json"
	"fmt"
	configuration "mongo-util/config"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetGridFsInfoByDB(config configuration.Mongo, dbName, collName string) (*[]configuration.GridFsAggregation, error) {

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.ConnectionString))
	if err != nil {
		panic(err)
	}
	defer client.Disconnect(ctx)

	db := client.Database(dbName)
	collNames := []string{collName}

	if names := getCollections(dbName, &collNames); names != nil && len(*names) > 0 {
		collNames = *names
	}

	var gridFsData []configuration.GridFsAggregation
	for _, collName := range collNames {
		collection := db.Collection(collName)
		data, err := getGridFsAggregationData(ctx, collection, collName)
		if err == nil || data != nil {
			for _, d := range *data {
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
	//var gridFsInfo []configuration.GridFsAggregation
	var gridFsInfo []bson.M
	if err = cursor.All(ctx, &gridFsInfo); err != nil {
		panic(err)
	}
	fmt.Println(gridFsInfo)
	return nil, err
}

func getCollections(dbName string, collNames *[]string) *[]string {
	return collNames
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
	//rslt, err := bson.Marshal(res)
	btRes, err := json.MarshalIndent(res, "", "	")

	var dbList configuration.DatabaseList
	//fmt.Println("b4 unmarshal:")
	if err := json.Unmarshal(btRes, &dbList); err != nil {
		fmt.Println("unmarshal err:", err)
		return nil, err
	}
	return &dbList.Databases, err
	//fmt.Println(len(dbList.Databases))
	//for _, db := range dbList.Databases {
	//
	//	database := client.Database(db.Name)
	//
	//	result, err := database.ListCollectionNames(
	//		context.TODO(),
	//		bson.D{})
	//
	//	if err != nil {
	//		log.Fatal(err)
	//	}

	//collRes, err := json.MarshalIndent(result, "", "	")
	//fmt.Println("dblist b4 unmarshal:", string(collRes))
	//
	//var collNames []interface{}
	//
	//if err := json.Unmarshal(btRes, &collNames); err != nil {
	//	fmt.Println("unmarshal err:", err)
	//	return
	//}
	//fmt.Printf("types : %+v", collNames)
	//for _, val := range collNames {
	//	episodesCollection := database.Collection(val.(string))
	//
	//	//id, _ := primitive.ObjectIDFromHex("5e3b37e51c9d4400004117e6")
	//
	//	//matchStage := bson.D{{"$match", bson.D{{"podcast", id}}}}
	//
	//	//groupStg := bson.D{"$group": bson.M{
	//	//	"_id":       "$contentType",
	//	//	"totalSize": bson.M{"sum": "$length"},
	//	//	"fileCount": bson.M{"sum": 1},
	//	//}}
	//	groupStage :=
	//		bson.D{
	//			{"$group", bson.D{
	//				{"_id", "$contentType"},
	//				{"totalSize", bson.D{{"$sum", "$length"}}},
	//				{"fileCount", bson.D{{"$sum", 1}}}}}}
	//	showInfoCursor, err := episodesCollection.Aggregate(ctx, mongo.Pipeline{groupStage})
	//	if err != nil {
	//		//panic(err)
	//	}
	//	//var showsWithInfo []bson.M
	//	var dbinfo DBInfo
	//	if err = showInfoCursor.All(ctx, &dbinfo); err != nil {
	//		//panic(err)
	//	}
	//		fmt.Printf(" dbinfo %+v", dbinfo)
	//	}
	//}

	//lookupStage := bson.D{{"$lookup", bson.D{{"from", "podcasts"}, {"localField", "podcast"}, {"foreignField", "_id"}, {"as", "podcast"}}}}
	//unwindStage := bson.D{{"$unwind", bson.D{{"path", "$podcast"}, {"preserveNullAndEmptyArrays", false}}}}
	//
	//showLoadedCursor, err := episodesCollection.Aggregate(ctx, mongo.Pipeline{lookupStage, unwindStage})
	//if err != nil {
	//	panic(err)
	//}
	//var showsLoaded []bson.M
	//if err = showLoadedCursor.All(ctx, &showsLoaded); err != nil {
	//	panic(err)
	//}
	//fmt.Println(showsLoaded)
	//
	//showLoadedStructCursor, err := episodesCollection.Aggregate(ctx, mongo.Pipeline{lookupStage, unwindStage})
	//if err != nil {
	//	panic(err)
	//}
	//var showsLoadedStruct []PodcastEpisode
	//if err = showLoadedStructCursor.All(ctx, &showsLoadedStruct); err != nil {
	//	panic(err)
	//}
	//fmt.Println(showsLoadedStruct)
}

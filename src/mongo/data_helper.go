package mongo

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	configuration "mongo-util/config"
	"net/http"
)

//This file not in use
func ReadAggregation(config configuration.Mongo) (interface{}, error) {
	//dataEndPoint := "https://data.mongodb-api.com/app/"
	url := fmt.Sprintf("https://data.mongodb-api.com/app/data-lopjk/endpoint/data/beta/action/aggregate")

	//aggrQuery := `db.fs.files.aggregate([   {$group:{_id:{DOMAIN_ID:"$DOMAIN_ID", contentType:"$contentType"},totalSize:{$sum:"$length"}, fileCount:{$sum: 1}}} ])`
	aggrQuery := `[{ $group :{_id:{_id:"$contentType"},totalSize:{$sum:"$length"},fileCount:{$sum:1}} }]`

	var queryType = configuration.Aggregation{
		DataSource: "zebra",
		Database:   "test",
		Collection: "uploads",
		Pipeline:   aggrQuery,
	}
	payload, _ := json.Marshal(queryType)
	log.Println("payload:", string(payload))
	//Make GET Call
	resp, err := configuration.HttpCall(http.MethodPost, url, payload, config)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == http.StatusOK {
		//log.Println(string(body))
		//err = json.Unmarshal(body, &data)
		log.Println(string(body))
		if err != nil {
			return nil, err
		}
	} else {
		if resp.StatusCode == http.StatusForbidden {
			err = errors.New("forbidden error, suggestion: check whether this machine IP is allowed to access the MongoCLuster")
			return nil, err
		}
		log.Fatalln("Mongo data API Error:", string(body))
	}
	return nil, err
}

func GetDatabaseInfo(config configuration.Mongo) (*[]configuration.Database, error) {
	return getDatabaseInfo(config)
}

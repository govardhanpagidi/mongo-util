package mongo

import (
	"errors"
	"fmt"
	"io/ioutil"
	configuration "mongo-util/config"
	"net/http"
	"strings"
)

func ExecuteQuery(config configuration.Mongo, query string) ([]byte, error) {
	url := fmt.Sprintf("%s/action/aggregate", config.DataEndPoint)
	method := "POST"

	payload := strings.NewReader(query)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Access-Control-Request-Headers", "*")
	req.Header.Add("api-key", config.ApiKey)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return body, errors.New(string(res.StatusCode))
	}

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	//fmt.Println(string(body))
	return body, err
}

func GetDatabaseInfo(config configuration.Mongo) (*[]configuration.Database, error) {
	dbs, err := getDatabaseInfo(config)
	if err != nil {
		return nil, err
	}
	var databases []configuration.Database
	for _, db := range *dbs {
		//Fetch list of collections in each DB
		if db.Name == nil {
			continue
		}
		collections, err := GetCollections(config, *db.Name)
		if err != nil {
			return nil, err
		}
		for _, coll := range *collections {
			db.Collections = append(db.Collections, coll.(string))
		}

		//Append collections to DB
		databases = append(databases, db)
	}
	return &databases, err
}

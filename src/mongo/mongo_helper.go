package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func getUsersByProject(config Mongo) ([]MongoUser, error) {
	url := fmt.Sprintf("%s/groups/%s/databaseUsers", config.AtlasEndPoint, config.ProjectID)

	var data UserData
	//Make GET Call
	resp, err := httpCall("GET", url, []byte(""), config)
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
		err = json.Unmarshal(body, &data)
		if err != nil {
			return nil, err
		}
	} else {
		if resp.StatusCode == http.StatusForbidden {
			err = errors.New("forbidden error, suggestion: check whether this machine IP is allowed to access the MongoCLuster")
			return nil, err
		}
		log.Fatalln("Mongo get users API Error:", string(body))
	}
	return data.Users, err
}

func getProjectByProjectName(config Mongo, projectName string) (*Project, error) {
	url := fmt.Sprintf("%s/groups/byName/%s", config.AtlasEndPoint, projectName)

	var project Project
	//Make GET Call
	resp, err := httpCall("GET", url, []byte(""), config)
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
		err = json.Unmarshal(body, &project)
		if err != nil {
			return nil, err
		}
	} else {
		if resp.StatusCode == http.StatusForbidden {
			err = errors.New("forbidden error, suggestion: check whether this machine IP is allowed to access the MongoCLuster")
			return nil, err
		}
		log.Fatalln("Mongo get users API Error:", string(body))
	}
	return &project, err
}

//Update mongo with new password for the user
func updatePassword(pwd string, user MongoUser, config Mongo) error {
	url := fmt.Sprintf("%s/groups/%s/databaseUsers/%s/%s", config.AtlasEndPoint, config.ProjectID, user.DBName, user.Username)

	//Generating payload for Atlas UpdateUser API
	data, err := json.Marshal(map[string]interface{}{
		"password": pwd,
	})
	if err != nil {
		log.Fatal("marshall error:", err)
	}

	//Make PATCH Call
	resp, err := httpCall(http.MethodPatch, url, data, config)
	if err != nil {
		log.Fatalln("error:", err)
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		//printing the body to get more error info
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		log.Println(string(body))
		//Return the status code rather error stack
		return errors.New(string(resp.StatusCode))
	}

	return err
}

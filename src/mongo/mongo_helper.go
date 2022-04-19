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

func GetUsersByProject(config configuration.Mongo) ([]configuration.MongoUser, error) {
	url := fmt.Sprintf("%s/groups/%s/databaseUsers", config.AtlasEndPoint, config.ProjectID)

	var data configuration.UserData
	//Make GET Call
	resp, err := configuration.HttpCall("GET", url, []byte(""), config)
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

func GetProjectByProjectName(config configuration.Mongo) (*configuration.Project, error) {
	url := fmt.Sprintf("%s/groups/byName/%s", config.AtlasEndPoint, config.ProjectName)

	var project configuration.Project
	//Make GET Call
	resp, err := configuration.HttpCall("GET", url, []byte(""), config)
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
		log.Fatalln("Mongo get Project Details by name API Error:", string(body))
	}
	return &project, err
}

//UpdatePassword is for updating db user password with random string
func UpdatePassword(pwd string, user configuration.MongoUser, config configuration.Mongo) error {
	url := fmt.Sprintf("%s/groups/%s/databaseUsers/%s/%s", config.AtlasEndPoint, config.ProjectID, user.DBName, user.Username)

	//Generating payload for Atlas UpdateUser API
	data, err := json.Marshal(map[string]interface{}{
		"password": pwd,
	})
	if err != nil {
		log.Fatal("marshall error:", err)
	}

	//Make PATCH Call
	resp, err := configuration.HttpCall(http.MethodPatch, url, data, config)
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

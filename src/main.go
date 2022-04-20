package main

import (
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

	args := os.Args[1:]
	if len(args) <= 0 || args[0] == "" {
		log.Fatalln("ProjectName is expected as an first argument!")
		return
	}
	config.Mongo.ProjectName = args[0]
}

func main() {
	//Read project name from args and get the project id
	project, err := mongo.GetProjectByProjectName(config.Mongo)
	if err != nil {
		log.Fatalln(err)
	}
	//log.Printf("Project Details %+v", project)

	//Setting projectId to config
	config.Mongo.ProjectID = project.ID
	//update with new password and push the same to GCP
	updateMongoUsers()
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

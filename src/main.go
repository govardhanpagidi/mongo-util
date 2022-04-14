package main

import (
	"log"
	"os"
)

var config Config

func init() {
	//This is for development purpose, config.json is expected in the same dir.
	if _, err := loadConfig("config.json"); err != nil {
		//Load configuration file, path can be changed according the file where it exists
		if _, err := loadConfig("/etc/config.json"); err != nil {
			log.Fatalln("config.json config load error:", err)
		}
	}

	//set gpc config
	err := os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", config.GCP.ConfigPath)
	if err != nil {
		log.Fatalln("GCP config load error:", err)
	}

	log.Println("config loaded successfully...")
}

func main() {
	//Read project name from args and get the project id
	args := os.Args[1:]
	if args[0] != "" {
		project, err := getProjectByProjectName(config.Mongo, args[0])
		if err != nil {
			log.Fatalln(err)
		}
		log.Printf("Project Details %+v", project)
		config.Mongo.ProjectID = project.ID
	}

	updateMongoUsers()
}

func updateMongoUsers() {
	//Fetch the list of mongodb users
	users, err := getUsersByProject(config.Mongo)
	if err != nil {
		log.Fatalln("get users:", err)
		return
	}
	log.Printf("Total users under %s : %d", config.Mongo.ProjectID, len(users))

	for _, userInfo := range users {

		pwd := RandomString(16)
		//log.Printf("Updating user %s with new password %s", result.Username, pwd)
		if err := updatePassword(pwd, userInfo, config.Mongo); err != nil {
			log.Printf("unable to change %s password for the DB: %s", userInfo.Username, userInfo.DBName)
			//log.Println(err)
			continue
		}

		//Save to GCP secret manager
		if err := saveSecret(config, userInfo, pwd); err != nil {
			log.Printf("gcp error: while updating the user %s for the DB %s :", userInfo.Username, userInfo.DBName)
			log.Println(err)
		}
		log.Printf("updated password for	 %s		%s", userInfo.Username, pwd)
	}
}

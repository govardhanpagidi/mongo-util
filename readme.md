# what is it meant for?
## A mongo db utility for
    
    1. updating the MongoDB users passwords and saving GCP Secret Manager
    2. Get MOngoDB cluster details as a report
    3. Get No.Of files and their sizes as a report
    4. Generate a report for a given aggreation query (WIP)

# Tech stack
    golang
    GCP Secret manager(Google Cloud Paltform) 
    MongoDB Atlas

# How does it work
    This is a mongodb utility to connect MongoDB and do some operations 
    using golang with the help of mongodb atlas api and mongodb drivers
    
    ./build-linux.sh with below params (if your os is linux flavoured) 
        -x <command> (Allowed values: "gridfs_report" "update_passwords" "cluster_report" "execute_query")
        -p <project_name> (Mongodb Atlas project Name)
        -d <database_name> (Mongo db Name)
        -c <cluster_name> (Atlas cluster name)
        -t <collection_name> 
        -k <data_api_key> (Atlas public key to connect data API)
        -b  <public_key> (Atlas public key to connect Atlas API)
        -r <private_key> (Atlas public key to connect Atlas API)
        -q <aggregation query> (Aggregation query in a valid Json format)
    Note: All parameters are not mandatory
        Below are the operation and their required parameters

###1) Update Passwords of DB users
    ./build-linux.sh -command update_passwords -p <project_name> 
        -p is mandatory
            Updates Passwords of all DB users of a given project
            Reads all the users in a given Project and updates the passwords if it has/can prevelige to change
            and update GCP secret manager 

            Note: currently it fails if you don't have GCP authentication credentilas in gcp.json

###2) Fetch reports with aggregation of ContentType, No.Of.Files and TotalSize
    ./build-linx.sh -command gridfs_report -d <databasename> -t <collectionname>   
        connection_string is mandatory in config.json (TODO : has to read from jenkins credentials)
        -d is optional (it considers all the db and their collections )
        -t is optional (it considers all the collections in a given db if you give database name only)

###3) Fetch reports with Cluster details
    ./build-linx.sh -command cluster_report -c <cluster name> -b <atlas public key>  -r <atlas private key>
        It uses Atlas API hence pub and priate keys of cluster is required
        
# How to Build and Run
##Prerequisites to build
### using jenkins?
    you need to setup
        1. Go runtime (Refer this if you want to setup in your local : [go installation] (https://go.dev/doc/install)
        2. If you want to integrate with jenkins use script.jenkinsfile and setup your job accordingly
            step 1 > manage jenkins > manage Plugins > search for "Go Plugin) and add it > Restart the Jenkins
            step 2 > Manage Jenkins > Global tool configuration > add go installation
            step 3 > You may need to add params manually according to the params used in scripts.jenkinsfile 
                        and just paster the content of scripts.jenkinsfile 
            step 4> you can see "build with parameters" once above steps are done, so just run it.
### Just run locally?
    you need to setup
       1. Go runtime (Refer this if you want to setup in your local : [go installation] (https://go.dev/doc/install)
            
            cd to ./src folder and just run 
            go build 
         
            Note: for the first time if it doesn't locate the dependents just run go mod tidy


# How to Test the functionality 
    You can follow the logs information for now :) 
        Check GCP secret manager is updating the passwords with versioning or not
        Check the Reports in Jenkins file path (Jenkins console logs will print where the source code is located so files stored in the same directory)

# Frequent issues
    1) You might face access related issues from your build machine to Atlas and Your key authoirzation to the mongodb resources
        a) you need to white list your build/jenkins machine Public IP in Atlas 
            (GO to Cloud.mongodb.com > select your organization > Select Project > Access Manager> 
                API keys > Edit your api key with Actions tab, dont for get click on Done button once you whitelist the IP)
        b) Your keys may not have project management access 
    2) Jenkins shouldn't have go plug-in installed as onetime setup
        step 1 > manage jenkins > manage Plugins > search for "Go Plugin) and add it > Restart the Jenkins
        step 2 > Manage Jenkins > Global tool configuration > add go installation
    3) Please check your jenkins credential ids to read all the credentials that you are referring, correct or not
    4) Your github repository structure is strictly followed in jenkins build scripts, please check once
    
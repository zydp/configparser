# configparser
	A tool used to read configuration files

# the configuration file format like this 
	#the server info
	[ServerInfo]
	listenHost=127.0.0.1
	listenPort=3244
	rootPath=/extern/url/getweburl

	#the database configuration
	[database]
	dbHost=xx.xx.xx.xx
	dbUser=test
	dbPasswd=123456
	dbName=testDB
	tableName=testTable

# use the configparser
	cfile, err := configparse.NewConfigInstance(configFileName)
	if nil != err {
		fmt.Println(err)
		return
	}
	serverHost := cfile.GetStrConfItem("ServerInfo", "listenHost")
	serverPort := cfile.GetIntegerConfItem("ServerInfo", "listenPort")

  
# if you want to judge whether it is right, you can use it like this
	serverHost := cfile.GetStrConfItem("ServerInfo", "listenHost")
	if "" == serverHost {
		fmt.Printf("No configuration items %s:%s\n", "ServerInfo", "listenHost" )
		return
	}
	serverPort := cfile.GetStrConfItem("ServerInfo", "listenHost")
	if serverHost<0 {
		fmt.Printf("No configuration items %s:%s\n", "ServerInfo", "listenHost" )
		return
	}

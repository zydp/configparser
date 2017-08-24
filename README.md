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

# example
	package main

	import (
		"fmt"
		"github.com/daipingpax/configparser"
	)

	func main() {
		cfile, err := configparser.NewConfigInstance("example.conf")
		if nil != err {
			fmt.Println(err)
			return
		}
		fmt.Println(cfile.GetStrConfItem("ServerInfo", "listenHost"))
		fmt.Println(cfile.GetIntegerConfItem("ServerInfo", "listenPort"))
		cfile.SetItemValue("NewModule", "TestStr", "a new item")
		cfile.SetItemValue("NewModule", "TestInt", 123)
		cfile.SetItemValue("NewModule", "TestFloat", 456.78)
		//cfile.SaveToFile("example.conf")		//auto rename the save file name is example.conf.new
		err = cfile.SaveToFile("newfile") //auto add suffix -> newfile.conf
		if nil != err {
			fmt.Println(err)
		}

		cfile.DelItem("database", "tableName")
		err = cfile.DelModule("NewModule")
		if nil != err {
			fmt.Println(err)
		}
		err = cfile.SaveToFile("newfile") //auto add suffix -> newfile.conf
		if nil != err {
			fmt.Println(err)
		}
	}

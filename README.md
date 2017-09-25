
# configparser
	一个简单配置文件读写模块,学习golang时用于练手。
	如有不适，欢迎提建议。谢谢！
	我的邮箱 daiping_zy@139.com
	

# in your command window

	go get github.com/zydp/configparser
	
# in your project

	import (
	    "github.com/zydp/configparser"	
	)
	
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

# using example
	package main

	import (
	    "fmt"
	    "github.com/zydp/configparser"
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

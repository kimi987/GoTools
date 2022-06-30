package conf

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

//Config conf
var Config struct {
	DBName      []string
	ConnAddr    []string //DB地址
	MaxOpenConn int
	MaxIdleConn int

	LogLevel string
	LogPath  string

	Usernames []string
	Passwords []string
}
var LoginUserMaps map[string]string

func init() {
	LoginUserMaps = make(map[string]string)
	data, err := ioutil.ReadFile("conf/config.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, &Config)
	if err != nil {
		panic(err)
	}
	for i := 0; i < len(Config.Usernames); i++ {
		fmt.Println("Init username = ", Config.Usernames[i])
		LoginUserMaps[Config.Usernames[i]] = Config.Passwords[i]
	}
}

//CheckLogin check can login
func CheckLogin(username, password string) bool {

	if LoginUserMaps[username] == "" {
		return false
	}
	if LoginUserMaps[username] == password {
		return true
	}
	return false
}

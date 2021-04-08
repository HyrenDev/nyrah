package utils

import (
	Databases "../../databases"
	"encoding/base64"
	"encoding/json"
	"github.com/gomodule/redigo/redis"
	"io/ioutil"
	"log"
	"os"
)

const (
	REDIS_ONLINE_COUNT_SCRIPT_PATH = "/home/scripts/redis/countOnlineUsers.lua"
)

func GetOnlinePlayers() int {
	script, err := ioutil.ReadFile(REDIS_ONLINE_COUNT_SCRIPT_PATH)

	if err != nil {
		panic(err)
	}

	redisConnection := Databases.StartRedis().Get()

	sha, err := redis.String(redisConnection.Do("SCRIPT", "LOAD", script))

	var redisData interface{}

	redisData, err = redisConnection.Do("EVALSHA", sha, "0")

	defer redisConnection.Close()

	if err != nil {
		log.Println("Couldn't get player count")

		return 0
	}

	return int(redisData.(int64))
}

func GetMaxPlayers() int {
	var settings = ReadSettingsFile()

	return int(settings["max_players"].(float64))
}

func GetFavicon() (string, error) {
	path, err := os.Getwd()

	if err != nil {
		log.Println(err)
	}

	file, err := ioutil.ReadFile(path + "/public/favicon.png")

	if err != nil {
		panic(err)
	}

	b64 := base64.StdEncoding.EncodeToString(file)
	output := "data:image/png;base64," + b64

	return output, nil
}

func GetServerAddress() string {
	var settings = ReadSettingsFile()

	return settings["address"].(string)
}

func GetServerPort() int {
	var settings = ReadSettingsFile()

	return int(settings["port"].(float64))
}

func GetServerMOTD() string {
	var settings = ReadSettingsFile()

	var motd = settings["motd"].(map[string]interface{})

	return motd["first_line"].(string) + "\n" + motd["second_line"].(string)
}

func ReadSettingsFile() map[string]interface{} {
	path, err := os.Getwd()

	if err != nil {
		log.Println(err)
	}

	dat, err := ioutil.ReadFile(path + "/settings.json")

	if err != nil {
		panic(err)
	}

	var data map[string]interface{}

	err = json.Unmarshal(dat, &data)

	return data
}

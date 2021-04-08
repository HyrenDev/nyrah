package utils

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	Databases "../../databases"
)

func GetOnlinePlayers() int {
	redisConnection := Databases.StartRedis().Get()

	var cursor = "0"

	for ok := true; ok; ok = (cursor != "0") {
		result, err := redisConnection.Do("SCAN", []string{
			cursor,
			"MATCH",
			"users:*",
		})

		if err != nil {
			log.Println(err)

			return 0
		}

		log.Println(result)
	}

	defer redisConnection.Close()

	return 0
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

	return settings["address"].(map[string]interface{})["host"].(string)
}

func GetServerPort() int {
	var settings = ReadSettingsFile()

	return int(settings["address"].(map[string]interface{})["port"].(float64))
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

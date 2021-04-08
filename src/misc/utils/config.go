package utils

import (
	"../constants"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"gominet/chat"
	"io/ioutil"
	"log"
	"os"

	Databases "../../databases"
)

func GetMOTD() chat.TextComponent {
	db := Databases.StartPostgres()

	row, err := db.Query("SELECT \"first_line\", \"second_line\" FROM \"motd\" LIMIT 1")

	if err == nil && row.Next() {
		var first_line string
		var second_line string

		err := row.Scan(&first_line, &second_line)

		if err == nil {
			return chat.TextComponent{
				Text: fmt.Sprintf(
					"%s\n%s",
					first_line,
					second_line,
				),
			}
		}
	}

	return chat.TextComponent{
		Text: fmt.Sprintf("%s - Nyrah", constants.SERVER_NAME),
		Component: chat.Component{
			Color: chat.Yellow,
		},
	}
}

func GetOnlinePlayers() int {
	redisConnection := Databases.StartRedis().Get()

	var onlinePlayers int

	cursor := 0

	for ok := true; ok; ok = cursor != 0 {
		result, err := redis.Values(redisConnection.Do("SCAN", cursor, "MATCH", "users:*"))

		if err != nil {
			log.Println(err)

			return 0
		}

		cursor, _ = redis.Int(result[0], nil)
		keys, _ := redis.Strings(result[1], nil)

		onlinePlayers += len(keys)
	}

	defer redisConnection.Close()

	return onlinePlayers
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
		log.Println(err)
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
		log.Println(err)
	}

	var data map[string]interface{}

	err = json.Unmarshal(dat, &data)

	return data
}

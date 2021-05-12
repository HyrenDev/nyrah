package utils

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/patrickmn/go-cache"
	"io/ioutil"
	"log"
	"net/hyren/nyrah/minecraft/chat"
	"os"
	"time"

	Databases "net/hyren/nyrah/databases"
)

var (
	CACHE = cache.New(cache.NoExpiration, 10*time.Second)

	WHITELISTED_GROUPS = []string{
		"MASTER",
		"DIRECTOR",
		"MANAGER",
		"ADMINISTRATOR",
		"MODERATOR",
		"HELPER",
	}
)

func GetMOTD() chat.TextComponent {
	db := Databases.StartMariaDB()

	var maintenance = IsMaintenanceModeEnabled()

	motd, found := CACHE.Get("motd")

	if !found {
		row, err := db.Query("SELECT `first_line`, `second_line` FROM `motd` LIMIT 1")

		if err == nil && row.Next() {
			var first_line string
			var second_line string

			_ = row.Scan(&first_line, &second_line)

			if maintenance == true {
				motd = chat.TextComponent{
					Text: fmt.Sprintf(
						"%s\n%s",
						first_line,
						"§cO servidor atualmente está em manutenção.",
					),
				}
			} else {
				motd = chat.TextComponent{
					Text: fmt.Sprintf(
						"%s\n%s",
						first_line,
						second_line,
					),
				}
			}

			defer row.Close()

			CACHE.Set("motd", motd, 15*time.Second)
		}
	}

	defer db.Close()

	return motd.(chat.TextComponent)
}

func IsMaintenanceModeEnabled() bool {
	db := Databases.StartMariaDB()

	current_state, found := CACHE.Get("maintenance")

	if !found {
		row, err := db.Query("SELECT `current_state` FROM `maintenance` WHERE `application_name`='nyrah';")

		if err == nil && row.Next() {
			_ = row.Scan(&current_state)

			CACHE.Set("maintenance", current_state, 1*time.Second)

			defer row.Close()
		}
	}

	defer db.Close()

	return current_state.(bool)
}

func IsGroupWhitelisted(group_name string) bool {
	for _, item := range WHITELISTED_GROUPS {
		if item == group_name {
			return true
		}
	}

	return false
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
	max_players, found := CACHE.Get("max_players")

	if !found {
		db := Databases.StartMariaDB()

		row, err := db.Query("SELECT `slots` FROM `applications` WHERE `name`='nyrah';")

		if err != nil {
			return 0
		}

		if row.Next() {
			_ = row.Scan(&max_players)
		}

		CACHE.Set("max_players", max_players, 3*time.Second)

		defer row.Close()
		defer db.Close()
	}

	return int(max_players.([]uint8)[0])
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

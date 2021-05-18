package utils

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/hyren/nyrah/cache/local"
	"net/hyren/nyrah/minecraft/chat"
	"os"
	"strconv"
	"time"

	NyrahProvider "net/hyren/nyrah/misc/providers"
)

func GetMOTD() chat.TextComponent {
	motd, found := local.CACHE.Get("motd")

	if !found {
		row, err := NyrahProvider.MARIA_DB_MAIN.Provide().Query("SELECT `first_line`, `second_line` FROM `motd` LIMIT 1")

		if err == nil && row.Next() {
			var firstLine string
			var secondLine string

			_ = row.Scan(&firstLine, &secondLine)

			var maintenance = IsMaintenanceModeEnabled()

			if maintenance == true {
				motd = chat.TextComponent{
					Text: fmt.Sprintf(
						"%s\n%s",
						firstLine,
						"§cO servidor atualmente está em manutenção.",
					),
				}
			} else {
				motd = chat.TextComponent{
					Text: fmt.Sprintf(
						"%s\n%s",
						firstLine,
						secondLine,
					),
				}
			}

			local.CACHE.Set("motd", motd, 5*time.Second)
		}
	}

	if motd == nil {
		return GetMOTD()
	}

	return motd.(chat.TextComponent)
}

func IsMaintenanceModeEnabled() bool {
	isMaintenanceModeEnabled, found := local.CACHE.Get("maintenance")

	if !found {
		row, err := NyrahProvider.MARIA_DB_MAIN.Provide().Query("SELECT `current_state` FROM `maintenance` WHERE `application_name`='nyrah';")

		if err == nil && row.Next() {
			var currentState bool

			row.Scan(&currentState)

			defer row.Close()

			isMaintenanceModeEnabled = currentState

			local.CACHE.Set("maintenance", isMaintenanceModeEnabled, 1*time.Second)
		}
	}

	if isMaintenanceModeEnabled == nil {
		return IsMaintenanceModeEnabled()
	}

	return isMaintenanceModeEnabled.(bool)
}

func GetMaxPlayers() int {
	maxPlayers, found := local.CACHE.Get("max_players")

	if !found {
		row, err := NyrahProvider.MARIA_DB_MAIN.Provide().Query("SELECT `slots` FROM `applications` WHERE `name`='nyrah';")

		if err != nil {
			return 0
		}

		if row.Next() {
			_ = row.Scan(&maxPlayers)
		}

		local.CACHE.Set("max_players", maxPlayers, 3*time.Second)

		defer row.Close()
	}

	maxPlayersValue, _ := strconv.Atoi(string(maxPlayers.([]byte)))

	return maxPlayersValue
}

func GetFavicon() (string, error) {
	serverIcon, found := local.CACHE.Get("server_icon")

	if !found {
		path, err := os.Getwd()

		if err != nil {
			log.Println(err)
		}

		file, err := ioutil.ReadFile(path + "/public/favicon.png")

		if err != nil {
			log.Println(err)
		}

		b64 := base64.StdEncoding.EncodeToString(file)
		serverIcon = "data:image/png;base64," + b64

		local.CACHE.Set("server_icon", serverIcon, 5*time.Minute)
	}

	return serverIcon.(string), nil
}

func GetServerAddress() string {
	var settings = readSettingsFile()

	return settings["inet_socket_address"].(map[string]interface{})["host"].(string)
}

func GetServerPort() int {
	var settings = readSettingsFile()

	return int(settings["inet_socket_address"].(map[string]interface{})["port"].(float64))
}

func readSettingsFile() map[string]interface{} {
	settings, found := local.CACHE.Get("settings")

	if !found {
		path, err := os.Getwd()

		if err != nil {
			log.Println(err)
		}

		file, err := os.ReadFile(path + "/settings.json")

		if err != nil {
			log.Println(err)
		}

		err = json.Unmarshal(file, &settings)

		local.CACHE.Set("settings", settings, 15*time.Hour)
	}

	return settings.(map[string]interface{})
}

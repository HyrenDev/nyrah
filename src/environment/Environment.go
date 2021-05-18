package environment

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/hyren/nyrah/cache/local"
	"time"
)

func getEnvironment() map[string]interface{} {
	environment, found := local.CACHE.Get("environment")

	if !found {
		bytes, err := ioutil.ReadFile("/home/configuration/environment.json")

		if err != nil {
			fmt.Println(err)
		}

		err = json.Unmarshal(bytes, &environment)

		local.CACHE.Set("environment", environment, 5*time.Minute)
	}

	if environment == nil {
		return getEnvironment()
	}

	return environment.(map[string]interface{})
}

func Get(key string) interface{} {
	return getEnvironment()[key]
}
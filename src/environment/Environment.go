package environment

import (
	"encoding/json"
	"io/ioutil"
	"net/hyren/nyrah/cache/local"
	"time"
)

var (
	environment map[string]interface{}
)

func loadEnvironment() map[string]interface{} {
	bytes, err := ioutil.ReadFile("/home/configuration/environment.json")

	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(bytes, &environment)

	local.CACHE.Set("environment", environment, 5*time.Minute)

	return environment
}

func Get(key string) interface{} {
	if environment == nil {
		loadEnvironment()
	}

	return environment[key]
}

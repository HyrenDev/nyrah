package environment

import (
	"encoding/json"
	"io/ioutil"
)

const (
	ENV_PATH = "/home/configuration/environment.json"
)

func ReadFile() map[string]interface{} {
	dat, err := ioutil.ReadFile(ENV_PATH)

	if err != nil {
		panic(err)
	}

	var data map[string]interface{}

	err = json.Unmarshal(dat, &data)

	return data
}
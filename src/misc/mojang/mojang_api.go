package mojang

import (
	_ "../utils"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	API_END_POINT = "https://api.minetools.eu"
)

func GetUniqueId(name string) map[string]interface{} {
	resp, err := http.Get(fmt.Sprintf(
		"%s/uuid/%s",
		API_END_POINT,
		name,
	))

	if err != nil {
		panic(err)
	}

	dat, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		panic(err)
	}

	var data map[string]interface{}

	err = json.Unmarshal(dat, &data)

	return data
}
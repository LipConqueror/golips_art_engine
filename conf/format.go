// format
package conf

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"golips_art_engine/models"
)

func GetConfig(debug bool) (*models.Config, error) {
	file, err := os.Open("./conf/config.json")

	if err != nil {
		if debug {
			log.Println("[OpenFile]", err)
		}
		return nil, err
	}

	defer file.Close()

	body, err := ioutil.ReadAll(file)

	if err != nil {
		if debug {
			log.Println("[ReadFile]", err)
		}
		return nil, err
	}

	var result models.Config

	err = json.Unmarshal(body, &result)

	if err != nil {
		if debug {
			log.Println("[ReadFile]", err)
		}
		return nil, err
	}

	return &result, nil
}

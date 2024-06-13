package vars

import (
	"encoding/json"
	"os"
	"strconv"

	"github.com/kercre123/wire-pod/chipper/pkg/logger"
)

var Language string = "en-US"

var IntentList []JsonIntent

type JsonIntent struct {
	Name              string   `json:"name"`
	Keyphrases        []string `json:"keyphrases"`
	RequireExactMatch bool     `json:"requiresexact"`
}

func LoadIntents() ([]JsonIntent, error) {
	jsonFile, err := os.ReadFile("intent-data/" + Language + ".json")
	var matches [][]string
	var jsonIntents []JsonIntent
	if err == nil {
		err = json.Unmarshal(jsonFile, &jsonIntents)
		if err != nil {
			logger.Println("Failed to load intents: " + err.Error())
		}

		for _, element := range jsonIntents {
			//logger.Println("Loading intent " + strconv.Itoa(index) + " --> " + element.Name + "( " + strconv.Itoa(len(element.Keyphrases)) + " keyphrases )")
			matches = append(matches, element.Keyphrases)
		}
		logger.Println("Loaded " + strconv.Itoa(len(jsonIntents)) + " intents and " + strconv.Itoa(len(matches)) + " matches (language: " + Language + ")")
	}
	return jsonIntents, err
}

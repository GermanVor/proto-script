package common

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

type Substitution struct {
	Substitution string `json:"substitution"`
	ImportSource string `json:"import_source"`
}

type SubstitutionMap = map[string]Substitution

type ConfigObj = map[string]SubstitutionMap

func ParseConfig(key string) (SubstitutionMap, bool) {
	jsonFile, err := os.Open("config.json")
	if err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()

	jsonData, err := io.ReadAll(jsonFile)
	if err != nil {
		log.Fatal(err)
	}

	config := ConfigObj{}

	if err := json.Unmarshal(jsonData, &config); err != nil {
		log.Fatal(err)
	}

	substitutionMap, ok := config[key]
	return substitutionMap, ok
}

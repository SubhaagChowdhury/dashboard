package config

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/go-playground/validator/v10"
)

// LoadCompleteConfigurations loads the complete configurations from a file.
func (configParserObject *Configurations) LoadCompleteConfigurations(configurationFilePath string) bool {
	configFileHandler, _ := os.Open(configurationFilePath)
	defer configFileHandler.Close()

	jsonParser := json.NewDecoder(configFileHandler)
	err := jsonParser.Decode(&configParserObject)
	if err != nil {
		pc, file, line, _ := runtime.Caller(0)
		log.Printf("%s in %s (%s:%d): %v\n", "Invalid data present in file", runtime.FuncForPC(pc).Name(), filepath.Base(file), line, err)
		return false
	}

	validate := validator.New()
	err = validate.Struct(configParserObject)
	if err != nil {
		pc, file, line, _ := runtime.Caller(0)
		log.Printf("%s in %s (%s:%d): %v\n", "Error validating struct", runtime.FuncForPC(pc).Name(), filepath.Base(file), line, err)
		return false
	}

	return true
}

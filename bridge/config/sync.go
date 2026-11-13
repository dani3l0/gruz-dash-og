package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// Saves in-memory config data to persistent storage
func Sync() {
	fmt.Println("Syncing in-memory config to disk ...")
	marshaled, _ := json.Marshal(Main)
	err := os.WriteFile(configFilename, marshaled, fileMode)
	if err != nil {
		fmt.Printf("Failed writing config file: %v\n", err)
	}
}

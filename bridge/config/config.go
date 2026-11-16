package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"
)

// Some constants
var configFilename = "config.json"
var fileMode = os.FileMode(0o750)
var Main = ConfigDefaults

// Read config or create a new config file with defaults
func Init() {
	f, err := os.ReadFile(configFilename)
	if os.IsNotExist(err) {
		fmt.Println(configFilename + " does not exist, creating one ...")
		jsoned, _ := json.Marshal(ConfigDefaults)
		err = os.WriteFile(configFilename, jsoned, fileMode)
		if err != nil {
			fmt.Println("Failed to create " + configFilename + ", exiting now ...")
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}
		fmt.Println("Config file has been created, have fun!")
	} else {
		err = json.Unmarshal(f, &Main)
		if err != nil {
			fmt.Println("Failed to read " + configFilename + ", exiting now ...")
			fmt.Printf("%v\n", err)
			os.Exit(0)
		}
		fmt.Println("Config file loaded, have fun!")
	}
	if Main.SerialPort == "auto" || Main.SerialPort == "" {
		Main.SerialPort = detectPort()
	}
}

// Automatically detects Serial port (easy-way)
func detectPort() string {
	entries, err := os.ReadDir("/dev/")
	if err != nil {
		fmt.Println("Failed to read /dev/: 'auto' setting for SerialPort is available only under Linux-based OS!")
		os.Exit(1)
	}
	for _, v := range entries {
		n := v.Name()
		if strings.HasPrefix(n, "ttyUSB") || strings.HasPrefix(n, "ttyACM") {
			fullpath := path.Join("/dev/", n)
			fmt.Println("Autodetected serial port: " + fullpath)
			return fullpath
		}
	}
	return ""
}

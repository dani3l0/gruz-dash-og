package config

// Config value types
type Config struct {
	Address    string
	Port       int
	SerialPort string
	Odometer   uint32
}

// Defaults configuration values
var ConfigDefaults = Config{
	Address:    "127.0.0.1",
	Port:       4444,
	SerialPort: "auto",
	Odometer:   0,
}

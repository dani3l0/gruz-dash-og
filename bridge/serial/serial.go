package serial

import (
	"fmt"

	"go.bug.st/serial"
)

var connection serial.Port
var connected bool

// Connects to target serial port
func Connect(target string, baudrate int) error {
	mode := &serial.Mode{
		BaudRate: baudrate,
		InitialStatusBits: &serial.ModemOutputBits{
			DTR: false,
			RTS: false,
		},
	}
	var err error
	connection, err = serial.Open(target, mode)
	if err != nil {
		fmt.Printf("Failed opening serial port (%s): %v\n", target, err)
	} else {
		connected = true
	}
	return err
}

// Writes data to connected serial port
func Write(data []byte) error {
	var err error
	if connected {
		connection.ResetInputBuffer()
		connection.ResetOutputBuffer()
		_, err = connection.Write(data)
	}
	return err
}

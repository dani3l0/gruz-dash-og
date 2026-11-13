package main

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"os"
	"os/signal"
	"outgouge/config"
	"outgouge/outgauge"
	"outgouge/serial"
	"syscall"
	"time"
)

func main() {
	// Load configuration
	config.Init()

	// Connect to UDP stream
	err := outgauge.Connect(config.Main.Address, config.Main.Port)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	} else {
		fmt.Printf("Waiting for UDP stream ...\n")
	}

	// Connect to Serial port
	err = serial.Connect(config.Main.SerialPort, 115200)
	if err != nil {
		fmt.Printf("Failed to connect to port, but continuing\n")
		time.Sleep(time.Second * 3)
	}

	// Handle CTRL+C or SIGTERM
	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	context.AfterFunc(ctx, func() {
		config.Sync()
		os.Exit(0)
	})

	// Variables needed for mileage calculation
	lastLoopTime := time.Now().UnixMicro()
	odoFloat := float64(config.Main.Odometer)

	// Main loop
	for {
		now := time.Now().UnixMicro()
		diff := now - lastLoopTime
		partOfTheSecond := float64(diff) / 1000.0 / 1000.0

		// Receive new data
		if err := outgauge.Receive(); err == nil {
			outgauge.LogData(outgauge.Parse(outgauge.RawData))
			p := outgauge.StripDownPacket(outgauge.RawData)
			buf := new(bytes.Buffer)
			binary.Write(buf, binary.LittleEndian, &p)
			byts := buf.Bytes()
			serial.Write(byts)
		}

		// If speed is more than 5cm/s, add driven distance to mileage
		if outgauge.RawData.Speed > 0.05 {
			odoFloat += partOfTheSecond * float64(outgauge.RawData.Speed)
		}
		config.Main.Odometer = uint32(odoFloat)
		time.Sleep(time.Millisecond)
		lastLoopTime = now
	}
}

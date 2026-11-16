package main

import (
	"context"
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
	lastPacketTime := time.Now().Add(-time.Minute)
	lastRealPacketTime := lastPacketTime
	odoFloat := float64(config.Main.Odometer)

	// OutGauge UDP stream loop
	go func() {
		for {
			// Receive new data and write it to serial
			if err := outgauge.Receive(); err == nil {
				outgauge.LogData(outgauge.Parse(outgauge.RawData))
				serial.Write(outgauge.ToSerial(outgauge.RawData))
				lastPacketTime = time.Now()
				lastRealPacketTime = lastPacketTime
			}
			time.Sleep(time.Millisecond)
		}
	}()

	// Main loop
	for {
		tnow := time.Now()
		now := tnow.UnixMicro()
		diff := now - lastLoopTime
		partOfTheSecond := float64(diff) / 1000.0 / 1000.0

		// If speed is more than 5cm/s and there is fresh packet, add driven distance to mileage
		if outgauge.RawData.Speed > 0.05 && tnow.Before(lastRealPacketTime.Add(time.Second)) {
			odoFloat += partOfTheSecond * float64(outgauge.RawData.Speed)
		}
		config.Main.Odometer = uint32(odoFloat)

		// If no packet has been sent for over 1 second, repeat
		if tnow.After(lastPacketTime.Add(time.Second * 2)) {
			serial.Write(outgauge.ToSerial(outgauge.RawData))
			lastPacketTime = tnow
		}

		// Finish loop
		time.Sleep(time.Millisecond * 5)
		lastLoopTime = now
	}
}

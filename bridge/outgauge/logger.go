package outgauge

import "fmt"

// Helper function to show dashboard icon
func dashLight(icon string, on bool) string {
	str := " "
	if on {
		str = icon
	}
	return str
}

// Pretty-prints OutGauge data to terminal
func LogData(data OutGaugeData) {
	// Clear terminal screen
	fmt.Print("\033[H\033[2J")

	// Print everything
	fmt.Printf("--- OUTGAUGE DASHBOARD ---\n")
	fmt.Printf("   %s ", dashLight("⬅️", data.SignalL))
	fmt.Printf(" %s ", dashLight("🪫", data.Battery))
	fmt.Printf(" %s ", dashLight("⚠️", data.Handbrake))
	fmt.Printf(" %s ", dashLight("🛑", data.ABS))
	fmt.Printf(" %s ", dashLight("🛞", data.TC))
	fmt.Printf(" %s ", dashLight("💡", data.FullBeam))
	fmt.Printf(" %s \n", dashLight("➡️", data.SignalR))
	fmt.Printf("Speed:         %.1f KM/H\n", data.Speed)
	fmt.Printf("RPM:           %d\n", data.RPM)
	fmt.Printf("Gear:          %d\n", data.Gear)
	fmt.Printf("Boost:         %.2f BAR\n", data.Boost)
	fmt.Printf("Fuel:          %.1f%%\n", data.Fuel*100)
	fmt.Printf("Engine Temp:   %.1f C\n", data.EngineTemp)
	fmt.Printf("Oil Temp:      %.1f C\n", data.OilTemp)
	fmt.Printf("Odometer:      %.3f KM\n", float32(data.Odometer)/1000)
}

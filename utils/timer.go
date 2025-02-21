package utils

import (
	"Assigment-1/config"
	"fmt"
	"time"
)

func StartUptime() {
	config.StartTime = time.Now()
}

func GetUptime() string {
	uptime := time.Since(config.StartTime)
	hours := int(uptime.Hours())
	minutes := int(uptime.Minutes()) % 60
	seconds := int(uptime.Seconds()) % 60

	return fmt.Sprintf("%02dh:%02dm:%02ds", hours, minutes, seconds)
}

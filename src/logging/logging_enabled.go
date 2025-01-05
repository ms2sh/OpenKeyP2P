//go:build logging
// +build logging

package logging

import "log"

// LogInfo protokolliert Informationsmeldungen
func LogInfo(format string, v ...interface{}) {
	log.Printf("[INFO] "+format, v...)
}

// LogError protokolliert Fehlermeldungen
func LogError(format string, v ...interface{}) {
	log.Printf("[ERROR] "+format, v...)
}

// LogDebug protokolliert Debug-Meldungen
func LogDebug(format string, v ...interface{}) {
	log.Printf("[DEBUG] "+format, v...)
}

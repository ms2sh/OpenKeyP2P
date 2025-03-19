package logging

import (
	"log"

	openkeyp2p "github.com/ms2sh/OpenKeyP2P/src"
)

// LogInfo protokolliert Informationsmeldungen
func LogInfo(llevel openkeyp2p.LogLevel, format string, v ...interface{}) {
	log.Printf("[INFO] "+format, v...)
}

// LogError protokolliert Fehlermeldungen
func LogError(llevel openkeyp2p.LogLevel, format string, v ...interface{}) {
	log.Printf("[ERROR] "+format, v...)
}

// LogDebug protokolliert Debug-Meldungen
func LogDebug(llevel openkeyp2p.LogLevel, format string, v ...interface{}) {
	log.Printf("[DEBUG] "+format, v...)
}

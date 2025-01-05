//go:build !logging
// +build !logging

package logging

// LogInfo ist eine No-Op Funktion, wenn Logging deaktiviert ist
func LogInfo(format string, v ...interface{}) {}

// LogError ist eine No-Op Funktion, wenn Logging deaktiviert ist
func LogError(format string, v ...interface{}) {}

// LogDebug ist eine No-Op Funktion, wenn Logging deaktiviert ist
func LogDebug(format string, v ...interface{}) {}

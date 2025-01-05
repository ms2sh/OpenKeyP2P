package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
)

// Lade TLS-Konfiguration mit Root CA für die Validierung der Client-Zertifikate
func loadTLSConfig(rootCAFile, certFile, keyFile string) (*tls.Config, error) {
	// Root CA Zertifikat laden
	rootCACert, err := os.ReadFile(rootCAFile)
	if err != nil {
		return nil, fmt.Errorf("Fehler beim Laden des Root CA Zertifikats: %w", err)
	}

	// Zertifikat und privater Schlüssel des Servers laden
	serverCert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, fmt.Errorf("Fehler beim Laden des Server-Zertifikats und Schlüssels: %w", err)
	}

	// Zertifikatspool für Root CA erstellen
	rootCAs := x509.NewCertPool()
	ok := rootCAs.AppendCertsFromPEM(rootCACert)
	if !ok {
		return nil, fmt.Errorf("Fehler beim Hinzufügen der Root CA zum Zertifikatspool")
	}

	// TLS-Konfiguration erstellen
	tlsConfig := &tls.Config{
		ClientCAs:          rootCAs,                        // Verwende Root CA für Client-Verifizierung
		ClientAuth:         tls.RequireAndVerifyClientCert, // Verlange Client-Zertifikat
		RootCAs:            rootCAs,                        // Verwende Root CA für Server-Authentifizierung
		InsecureSkipVerify: false,
		Certificates:       []tls.Certificate{serverCert}, // Server Zertifikat
	}

	return tlsConfig, nil
}

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/libp2p/go-libp2p-core/crypto"
)

type Config struct {
	Peers []string `json:"peers"`
}

func loadConfig(filePath string) (*Config, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("Fehler beim Lesen der Datei: %w", err)
	}
	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("Fehler beim Parsen der Konfiguration: %w", err)
	}
	return &config, nil
}

func loadOrGenerateKey(path string) (crypto.PrivKey, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		priv, _, err := crypto.GenerateEd25519Key(nil)
		if err != nil {
			return nil, fmt.Errorf("Fehler beim Generieren des Schlüssels: %w", err)
		}
		data, _ := crypto.MarshalPrivateKey(priv)
		_ = ioutil.WriteFile(path, data, 0600)
		return priv, nil
	}
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("Fehler beim Laden des Schlüssels: %w", err)
	}
	return crypto.UnmarshalPrivateKey(data)
}

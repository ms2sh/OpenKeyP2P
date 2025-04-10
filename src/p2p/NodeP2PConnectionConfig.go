package p2p

import (
	"fmt"
	"net/url"
	"slices"
	"strings"
)

// Erstellt aus 2 Peer Configurationen 1ne
func _DeterminesCommonConfig(remoteControlStream NodeP2PConnectionConfig, localNodeConfig NodeP2PConnectionConfig) NodeP2PConnectionConfig {
	// Die Listen von Einträgen aus beiden Konfigurationen holen
	remoteList := remoteControlStream.List()
	localList := localNodeConfig.List()

	// Eine Map erstellen, um die gemeinsamen Einträge zu speichern
	commonEntries := NewNodeP2PConnectionConfig()

	// Eine Hilfsfunktion, die überprüft, ob ein Eintrag in einer Liste enthalten ist
	contains := func(entry NodeP2PConfigEntry, list []NodeP2PConfigEntry) bool {
		return slices.Contains(list, entry)
	}

	// Iteriere durch die Einträge der remote-Konfiguration und finde gemeinsame Einträge
	for _, remoteEntry := range remoteList {
		if contains(remoteEntry, localList) {
			// Wenn der Eintrag auch in der lokalen Liste vorhanden ist, füge ihn zu den gemeinsamen Einträgen hinzu
			commonEntries.Add(remoteEntry.Name, remoteEntry.Value)
		}
	}

	return commonEntries
}

// Erstellt eine neue Konfiguration für eine Verbindung
func NewNodeP2PConnectionConfig() NodeP2PConnectionConfig {
	return NodeP2PConnectionConfig("<>")
}

// Fügt ein neues name=value Paar innerhalb der <> hinzu (ähnlich wie Cookies)
func (o *NodeP2PConnectionConfig) Add(name string, value string) error {
	// Überprüfe, ob der Name gültig ist
	if !isValidName(name) {
		return fmt.Errorf("invalid name: '%s' only a-z, A-Z, 0-9, _ and - are allowed", name)
	}

	// URL-codiert den Wert, um Sonderzeichen zu behandeln
	encodedValue := url.QueryEscape(value)

	// Erstelle den neuen Eintrag im Format "name=value"
	entry := fmt.Sprintf("%s=%s", name, encodedValue)

	// Falls die Konfiguration noch leer ist, füge den Eintrag zwischen <>
	if string(*o) == "<>" {
		*o = NodeP2PConnectionConfig("<" + entry + ">")
	} else {
		// Falls bereits Einträge vorhanden sind, füge den neuen Eintrag hinzu,
		// und trenne sie mit einem Semikolon (ähnlich wie bei Cookies)
		// Dereferenzieren von o und den String bearbeiten
		trimmedConfig := string(*o)[:len(string(*o))-1] // Entfernt das schließende ">"
		*o = NodeP2PConnectionConfig(fmt.Sprintf("%s;%s>", trimmedConfig, entry))
	}

	return nil
}

// Gibt den Wert für den angegebenen Namen zurück, falls vorhanden
func (o NodeP2PConnectionConfig) Get(name string) string {
	// Entferne das öffnende "<" und das schließende ">"
	trimmedConfig := string(o)[1 : len(string(o))-1]

	// Durchsuche die Einträge nach dem gewünschten Namen
	pairs := strings.Split(trimmedConfig, ";")
	for _, pair := range pairs {
		parts := strings.Split(pair, "=")
		if len(parts) == 2 && parts[0] == name {
			// Dekodiere den Wert, bevor er zurückgegeben wird
			decodedValue, err := url.QueryUnescape(parts[1])
			if err != nil {
				return ""
			}
			return decodedValue
		}
	}

	// Wenn kein Wert gefunden wurde, gib einen leeren String zurück
	return ""
}

// Listet alle Einträge als ConfigEntry Structs auf
func (o NodeP2PConnectionConfig) List() []NodeP2PConfigEntry {
	// Entferne das öffnende "<" und das schließende ">"
	trimmedConfig := string(o)[1 : len(string(o))-1]

	// Zerlege die Konfiguration in Paare (name=value)
	pairs := strings.Split(trimmedConfig, ";")
	var result []NodeP2PConfigEntry

	// Gehe durch die Paare und füge sie der Ergebnisliste hinzu
	for _, pair := range pairs {
		parts := strings.Split(pair, "=")
		if len(parts) == 2 {
			// Dekodiere den Wert
			decodedValue, err := url.QueryUnescape(parts[1])
			if err == nil {
				// Erstelle das ConfigEntry und füge es der Liste hinzu
				result = append(result, NodeP2PConfigEntry{
					Name:  parts[0],
					Value: decodedValue,
				})
			}
		}
	}

	return result
}

// Wird verwendet um zu überprüfen ob eine Gewisse Einstellung vorhanden ist
func (o *NodeP2PConnectionConfig) HasConfigEntry(name string) bool {
	for _, item := range o.List() {
		if item.Name == name {
			return true
		}
	}
	return false
}

// Wird verwendet um zu überprüfen ob eine Gewisse Einstellung sammt Wert vorhanden ist
func (o *NodeP2PConnectionConfig) HasConfigEntryWithValue(name string, value string) bool {
	for _, item := range o.List() {
		if item.Name == name {
			if item.Value == value {
				return true
			} else {
				return false
			}
		}
	}
	return false
}

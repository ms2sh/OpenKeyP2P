package p2p

import (
	"encoding/json"
	"fmt"

	"github.com/ms2sh/OpenKeyP2P/src/logging"
)

// sendVersionInfo sends the version information over the connection.
func (c *FWConn) sendVersionInfo(v *VersionInfo) error {
	data, err := json.Marshal(v)
	if err != nil {
		logging.LogError("Error serializing version information: %v", err)
		return err
	}

	// Use the existing Write method to send the data.
	if err := c.Write(data); err != nil {
		logging.LogError("Error sending version information: %v", err)
		return err
	}

	if v.AgreedVersion != "" {
		logging.LogInfo("Agreed version sent: %s", v.AgreedVersion)
	} else {
		logging.LogInfo("Supported versions sent: %v", v.SupportedVersions)
	}
	return nil
}

// receiveClientVersionInfo receives the supported versions from the client.
func (c *FWConn) receiveClientVersionInfo(v *VersionInfo) error {
	data, err := c.Read()
	if err != nil {
		logging.LogError("Error reading client version information: %v", err)
		return err
	}

	if err := json.Unmarshal(data, v); err != nil {
		logging.LogError("Error deserializing client version information: %v", err)
		return err
	}

	logging.LogInfo("Received supported versions from client: %v", v.SupportedVersions)
	return nil
}

// receiveServerResponse receives the agreed version from the server.
func (c *FWConn) receiveServerResponse(v *VersionInfo) error {
	responseData, err := c.Read()
	if err != nil {
		logging.LogError("Error reading server response: %v", err)
		return err
	}

	var response VersionInfo
	if err := json.Unmarshal(responseData, &response); err != nil {
		logging.LogError("Error deserializing server response: %v", err)
		return err
	}

	if response.AgreedVersion == "" {
		logging.LogError("No agreed version received from server")
		return fmt.Errorf("no agreed version received from server")
	}

	v.AgreedVersion = response.AgreedVersion
	return nil
}

package integration

import (
	"encoding/json"
	"fmt"
)

// GenerateDriverJSON builds a throwaway Integration, calls newClient on it, and marshals the
// DriverMetadata that call set. newClient is a driver's NewXClient constructor (e.g.
// deconzclient.NewDeconzClient) wrapped to discard its return value - the same call every driver
// binary makes at startup, which only builds DriverMetadata and registers function pointers, so
// it's safe to call here without InitClient/Run.
//
// The result is the driver.json a custom-installed driver archive needs at its root (see
// doc/integration-driver/driver-installation.md). Generating it from the exact same constructor a
// running driver uses means the two can never drift apart.
func GenerateDriverJSON(newClient func(*Integration)) ([]byte, error) {
	i, err := NewIntegration(Config{})
	if err != nil {
		return nil, fmt.Errorf("NewIntegration: %w", err)
	}

	newClient(i)

	if i.Metadata == nil {
		return nil, fmt.Errorf("client constructor didn't set driver metadata")
	}

	data, err := json.MarshalIndent(i.Metadata, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal driver metadata: %w", err)
	}
	data = append(data, '\n')

	return data, nil
}

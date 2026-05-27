package bridge

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const seederRegistryURL = "https://gist.githubusercontent.com/dipghoshraj/dbb415d4987a99ef465add5945cca071/raw/config.json"

type SeederRegion struct {
	Country string `json:"country"`
	Flag    string `json:"flag"`
}

type SeederInfo struct {
	IP          string       `json:"ip"`
	Fingerprint string       `json:"fingerprint"`
	Status      string       `json:"status"`
	Maintainer  string       `json:"maintainer"`
	Region      SeederRegion `json:"region"`
}

type SeederInfoResponse struct {
	Network string       `json:"network"`
	Sources []string     `json:"sources"`
	Seeders []SeederInfo `json:"seeders"`
}

// ScanForSeeders fetches the list of registered Seeders from the AgniStack registry.
func ScanForSeeders() (SeederInfoResponse, error) {
	resp, err := http.Get(seederRegistryURL)
	if err != nil {
		return SeederInfoResponse{}, fmt.Errorf("failed to fetch seeder registry: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return SeederInfoResponse{}, fmt.Errorf("registry returned unexpected status: %d", resp.StatusCode)
	}

	var result SeederInfoResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return SeederInfoResponse{}, fmt.Errorf("failed to decode seeder registry response: %w", err)
	}

	return result, nil
}

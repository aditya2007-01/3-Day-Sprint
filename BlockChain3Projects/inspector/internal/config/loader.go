package config

import (
    "encoding/json"
    "fmt"
    "os"
)

type NetworkConfig struct {
    Nodes []NodeConfig `json:"nodes"`
}

type NodeConfig struct {
    Name   string `json:"name"`
    DBPath string `json:"db_path"`
    RPCURL string `json:"rpc_url,omitempty"`
}

func LoadConfig(path string) (*NetworkConfig, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, fmt.Errorf("failed to read config: %w", err)
    }

    var config NetworkConfig
    if err := json.Unmarshal(data, &config); err != nil {
        return nil, fmt.Errorf("failed to parse config: %w", err)
    }

    if len(config.Nodes) == 0 {
        return nil, fmt.Errorf("no nodes defined in config")
    }

    return &config, nil
}

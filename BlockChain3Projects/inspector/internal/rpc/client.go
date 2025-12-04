package rpc

import (
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "time"

    "inspector/internal/blocks"
)

type Client struct {
    baseURL string
    client  *http.Client
}

func NewClient(baseURL string) *Client {
    return &Client{
        baseURL: baseURL,
        client: &http.Client{
            Timeout: 10 * time.Second,
        },
    }
}

func (c *Client) FetchBlock(height int) (*blocks.Block, error) {
    url := fmt.Sprintf("%s/block/%d", c.baseURL, height)
    
    resp, err := c.client.Get(url)
    if err != nil {
        return nil, fmt.Errorf("failed to fetch block: %w", err)
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        body, _ := io.ReadAll(resp.Body)
        return nil, fmt.Errorf("RPC error: status %d, body: %s", resp.StatusCode, string(body))
    }
    
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("failed to read response: %w", err)
    }
    
    var block blocks.Block
    if err := json.Unmarshal(body, &block); err != nil {
        return nil, fmt.Errorf("failed to parse block: %w", err)
    }
    
    return &block, nil
}

type HealthResponse struct {
    Height        int     `json:"height"`
    LastBlockTime int64   `json:"last_block_time"`
    Peers         int     `json:"peers"`
    BlocksPerMin  float64 `json:"blocks_per_min"`
    Status        string  `json:"status"`
}

func (c *Client) FetchHealth() (*HealthResponse, error) {
    url := fmt.Sprintf("%s/health", c.baseURL)
    
    resp, err := c.client.Get(url)
    if err != nil {
        return nil, fmt.Errorf("failed to fetch health: %w", err)
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("health endpoint returned status %d", resp.StatusCode)
    }
    
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }
    
    var health HealthResponse
    if err := json.Unmarshal(body, &health); err != nil {
        return nil, err
    }
    
    return &health, nil
}

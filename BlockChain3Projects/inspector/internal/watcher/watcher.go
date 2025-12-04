package watcher

import (
    "fmt"
    "time"
    
    "inspector/internal/rpc"
)

const (
    ColorReset  = "\033[0m"
    ColorGreen  = "\033[32m"
    ColorYellow = "\033[33m"
    ColorRed    = "\033[31m"
    ColorCyan   = "\033[36m"
)

func Watch(rpcURL string, interval int) {
    client := rpc.NewClient(rpcURL)
    
    fmt.Printf("%s╔════════════════════════════════════════════════════════════════╗%s\n", ColorCyan, ColorReset)
    fmt.Printf("%s║            BHIV BLOCKCHAIN NODE WATCHER                       ║%s\n", ColorCyan, ColorReset)
    fmt.Printf("%s╚════════════════════════════════════════════════════════════════╝%s\n", ColorCyan, ColorReset)
    fmt.Printf("\nWatching: %s\n", rpcURL)
    fmt.Printf("Interval: %ds\n", interval)
    fmt.Println("Press Ctrl+C to stop\n")
    
    ticker := time.NewTicker(time.Duration(interval) * time.Second)
    defer ticker.Stop()
    
    fetchAndDisplay(client)
    
    for {
        select {
        case <-ticker.C:
            fetchAndDisplay(client)
        }
    }
}

func fetchAndDisplay(client *rpc.Client) {
    health, err := client.FetchHealth()
    if err != nil {
        fmt.Printf("%s[ERROR] Failed to fetch health: %v%s\n", ColorRed, err, ColorReset)
        return
    }
    
    displayHealth(health)
}

func displayHealth(health *rpc.HealthResponse) {
    timestamp := time.Now().Format("15:04:05")
    
    timeSinceBlock := time.Now().Unix() - health.LastBlockTime
    
    status := ColorGreen
    statusText := "HEALTHY"
    
    if timeSinceBlock > 60 {
        status = ColorRed
        statusText = "STUCK  "
    } else if timeSinceBlock > 30 {
        status = ColorYellow
        statusText = "SLOW   "
    }
    
    lastBlockTimeStr := time.Unix(health.LastBlockTime, 0).Format("15:04:05")
    
    fmt.Printf("[%s] %s%s%s | Height: %4d | Peers: %d | Last: %s (%ds ago) | Rate: %.1f blk/min\n",
        timestamp, 
        status, 
        statusText, 
        ColorReset, 
        health.Height, 
        health.Peers, 
        lastBlockTimeStr,
        timeSinceBlock,
        health.BlocksPerMin)
}

package consensus

import (
    "encoding/json"
    "fmt"
    "strings"
)

func OutputConsensusResult(result *ConsensusResult, jsonMode bool) {
    if jsonMode {
        jsonData, _ := json.MarshalIndent(result, "", "  ")
        fmt.Println(string(jsonData))
    } else {
        outputConsensusText(result)
    }
}

func outputConsensusText(result *ConsensusResult) {
    fmt.Println("\n" + strings.Repeat("═", 72))
    fmt.Println("CONSENSUS ANALYSIS REPORT")
    fmt.Println(strings.Repeat("═", 72))
    
    fmt.Printf("\n📊 NETWORK OVERVIEW:\n")
    fmt.Printf("  Scan Time:         %s\n", result.Timestamp)
    fmt.Printf("  Total Nodes:       %d\n", result.TotalNodes)
    fmt.Printf("  Network Health:    %s\n", getHealthEmoji(result.NetworkHealth))
    fmt.Printf("  Canonical Chain:   %s\n", result.CanonicalChain)
    fmt.Printf("  Consensus Height:  %d\n", result.ConsensusHeight)
    fmt.Printf("  Fork Points:       %d\n", len(result.ForkPoints))
    
    if len(result.ForkPoints) > 0 {
        fmt.Println("\n🔀 FORK ANALYSIS:")
        for i, fork := range result.ForkPoints {
            fmt.Printf("  Fork #%d:\n", i+1)
            fmt.Printf("    Height:   %d\n", fork.Height)
            fmt.Printf("    Branches: %d\n", fork.Branches)
            fmt.Printf("    Affected: %v\n", fork.AffectedNodes)
        }
    }
    
    fmt.Println("\n🖥️  NODE STATUS:")
    for name, state := range result.NodeStates {
        statusIcon := getNodeStatusIcon(state)
        
        fmt.Printf("  %s %s:\n", statusIcon, name)
        fmt.Printf("      Height:        %d\n", state.Height)
        fmt.Printf("      Status:        %s\n", state.Status)
        fmt.Printf("      Blocks Behind: %d\n", state.BlocksBehind)
        fmt.Printf("      On Canonical:  %v\n", state.OnCanonical)
    }
    
    fmt.Println("\n💡 RECOMMENDATIONS:")
    if len(result.Recommendations) == 0 {
        fmt.Println("  No issues detected - network is healthy")
    } else {
        for i, rec := range result.Recommendations {
            fmt.Printf("  %d. %s\n", i+1, rec)
        }
    }
    
    fmt.Println(strings.Repeat("═", 72))
}

func getHealthEmoji(health string) string {
    healthMap := map[string]string{
        "EXCELLENT": "✅ EXCELLENT",
        "GOOD":      "🟢 GOOD",
        "FAIR":      "🟡 FAIR",
        "POOR":      "🟠 POOR",
        "WARNING":   "⚠️  WARNING",
        "CRITICAL":  "🔴 CRITICAL",
    }
    
    if emoji, ok := healthMap[health]; ok {
        return emoji
    }
    return health
}

func getNodeStatusIcon(state NodeState) string {
    if !state.OnCanonical {
        return "❌"
    } else if state.BlocksBehind > 10 {
        return "🔴"
    } else if state.BlocksBehind > 0 {
        return "⚠️"
    } else {
        return "✅"
    }
}

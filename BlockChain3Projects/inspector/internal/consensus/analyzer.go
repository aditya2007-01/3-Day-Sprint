package consensus

import (
    "fmt"
    "time"

    "inspector/internal/db"
)

type NodeInfo struct {
    Name      string
    DBPath    string
    Height    int
    Storage   *db.Storage
}

type ConsensusResult struct {
    Timestamp           string              `json:"timestamp"`
    TotalNodes          int                 `json:"total_nodes"`
    CanonicalChain      string              `json:"canonical_chain"`
    ConsensusHeight     int                 `json:"consensus_height"`
    ForkPoints          []ForkPoint         `json:"fork_points"`
    NodeStates          map[string]NodeState `json:"node_states"`
    Recommendations     []string            `json:"recommendations"`
    NetworkHealth       string              `json:"network_health"`
}

type ForkPoint struct {
    Height        int      `json:"height"`
    Branches      int      `json:"branches"`
    AffectedNodes []string `json:"affected_nodes"`
}

type NodeState struct {
    Height       int      `json:"height"`
    Status       string   `json:"status"`
    BlocksBehind int      `json:"blocks_behind"`
    OnCanonical  bool     `json:"on_canonical"`
}

func AnalyzeConsensus(nodes []NodeInfo) (*ConsensusResult, error) {
    result := &ConsensusResult{
        Timestamp:  time.Now().Format("2006-01-02 15:04:05"),
        TotalNodes: len(nodes),
        NodeStates: make(map[string]NodeState),
        ForkPoints: []ForkPoint{},
    }

    if len(nodes) == 0 {
        return result, fmt.Errorf("no nodes provided")
    }

    maxHeight := 0
    for _, node := range nodes {
        if node.Height > maxHeight {
            maxHeight = node.Height
        }
    }
    
    consensusMap := make(map[int]map[string][]string)
    
    for height := 0; height <= maxHeight; height++ {
        consensusMap[height] = make(map[string][]string)
        
        for _, node := range nodes {
            if height <= node.Height {
                block, err := node.Storage.LoadBlock(height)
                if err == nil {
                    consensusMap[height][block.Hash] = append(
                        consensusMap[height][block.Hash],
                        node.Name,
                    )
                }
            }
        }
    }

    for height := 0; height <= maxHeight; height++ {
        if len(consensusMap[height]) > 1 {
            affectedNodes := []string{}
            for _, nodeList := range consensusMap[height] {
                affectedNodes = append(affectedNodes, nodeList...)
            }
            
            result.ForkPoints = append(result.ForkPoints, ForkPoint{
                Height:        height,
                Branches:      len(consensusMap[height]),
                AffectedNodes: affectedNodes,
            })
        }
    }

    result.CanonicalChain = findCanonicalChain(nodes, consensusMap)
    result.ConsensusHeight = findConsensusHeight(consensusMap, maxHeight)

    for _, node := range nodes {
        state := analyzeNodeState(node, result.CanonicalChain, maxHeight)
        result.NodeStates[node.Name] = state
    }

    result.Recommendations = generateConsensusRecommendations(result)
    result.NetworkHealth = calculateNetworkHealth(result)

    return result, nil
}

func findCanonicalChain(nodes []NodeInfo, consensusMap map[int]map[string][]string) string {
    maxScore := 0
    canonical := ""
    
    for _, node := range nodes {
        score := 0
        for height := 0; height <= node.Height; height++ {
            block, err := node.Storage.LoadBlock(height)
            if err == nil {
                score += len(consensusMap[height][block.Hash])
            }
        }
        
        if score > maxScore {
            maxScore = score
            canonical = node.Name
        }
    }
    
    return canonical
}

func findConsensusHeight(consensusMap map[int]map[string][]string, maxHeight int) int {
    for height := maxHeight; height >= 0; height-- {
        if len(consensusMap[height]) == 1 {
            return height
        }
    }
    return 0
}

func analyzeNodeState(node NodeInfo, canonical string, maxHeight int) NodeState {
    state := NodeState{
        Height:       node.Height,
        BlocksBehind: maxHeight - node.Height,
        OnCanonical:  (node.Name == canonical),
    }

    if node.Height == maxHeight {
        state.Status = "synchronized"
    } else if node.Height < maxHeight {
        state.Status = "behind"
    } else {
        state.Status = "ahead"
    }

    return state
}

func generateConsensusRecommendations(result *ConsensusResult) []string {
    recs := []string{}

    if len(result.ForkPoints) > 0 {
        recs = append(recs, fmt.Sprintf("⚠️  Fork detected at %d point(s)", len(result.ForkPoints)))
        for _, fork := range result.ForkPoints {
            recs = append(recs, fmt.Sprintf("   Height %d: %d branches affecting %d nodes", 
                fork.Height, fork.Branches, len(fork.AffectedNodes)))
        }
    }

    for name, state := range result.NodeStates {
        if !state.OnCanonical {
            recs = append(recs, fmt.Sprintf("🔧 %s: Resync from canonical chain (%s)", name, result.CanonicalChain))
        }
        if state.BlocksBehind > 10 {
            recs = append(recs, fmt.Sprintf("📥 %s: Critically behind - sync %d blocks urgently", name, state.BlocksBehind))
        } else if state.BlocksBehind > 0 {
            recs = append(recs, fmt.Sprintf("📥 %s: Sync %d blocks from network", name, state.BlocksBehind))
        }
    }

    if len(recs) == 0 {
        recs = append(recs, "✅ All nodes in perfect consensus - no action needed")
    }

    return recs
}

func calculateNetworkHealth(result *ConsensusResult) string {
    if len(result.ForkPoints) > 3 {
        return "CRITICAL"
    } else if len(result.ForkPoints) > 0 {
        return "WARNING"
    }
    
    syncedNodes := 0
    for _, state := range result.NodeStates {
        if state.BlocksBehind == 0 && state.OnCanonical {
            syncedNodes++
        }
    }
    
    syncPercentage := float64(syncedNodes) / float64(result.TotalNodes) * 100
    
    if syncPercentage >= 90 {
        return "EXCELLENT"
    } else if syncPercentage >= 70 {
        return "GOOD"
    } else if syncPercentage >= 50 {
        return "FAIR"
    } else {
        return "POOR"
    }
}

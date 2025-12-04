package main

import (
    "encoding/json"
    "flag"
    "fmt"
    "os"
    "time"

    "inspector/internal/blocks"
    "inspector/internal/config"
    "inspector/internal/consensus"
    "inspector/internal/db"
    "inspector/internal/errors"
    "inspector/internal/report"
    "inspector/internal/rpc"
    "inspector/internal/watcher"
)

const version = "1.0.0"

func main() {
    // Day 1 flags
    dbPath := flag.String("db", "./leveldb-data", "Path to LevelDB database")
    db1Path := flag.String("db1", "./node1-data", "Path to first database")
    db2Path := flag.String("db2", "./node2-data", "Path to second database")
    cmd := flag.String("cmd", "help", "Command: load, block, scan-errors, compare, consensus, watch, report")
    numBlocks := flag.Int("blocks", 10, "Number of blocks to load")
    jsonOutput := flag.Bool("json", false, "Output in JSON format")
    showVersion := flag.Bool("version", false, "Show version")
    
    // RPC flags
    rpcURL := flag.String("rpc", "", "RPC endpoint URL (e.g., http://localhost:8545)")
    watchInterval := flag.Int("interval", 2, "Watch mode interval in seconds")
    
    // Day 2 flags
    configPath := flag.String("config", "nodes.json", "Path to network config file")
    
    // Day 3 flags
    reportPath := flag.String("report", "inspector-report.json", "Output path for generated report")
    
    flag.Parse()

    if *showVersion {
        fmt.Printf("BHIV Chain Inspector v%s\n", version)
        return
    }

    switch *cmd {
    case "load":
        loadSampleData(*dbPath, *numBlocks)
    case "block":
        viewBlock(*dbPath, *rpcURL, *jsonOutput)
    case "scan-errors":
        runScan(*dbPath, *jsonOutput)
    case "compare":
        runCompare(*db1Path, *db2Path, *jsonOutput)
    case "consensus":
        runConsensus(*configPath, *jsonOutput)
    case "watch":
        runWatch(*rpcURL, *watchInterval)
    case "report":
        runFullReport(*configPath, *reportPath)
    case "help":
        printUsage()
    default:
        fmt.Printf("Unknown command: %s\n", *cmd)
        printUsage()
    }
}

func loadSampleData(dbPath string, numBlocks int) {
    storage, err := db.NewStorage(dbPath)
    if err != nil {
        fmt.Printf("❌ Error: %v\n", err)
        os.Exit(1)
    }
    defer storage.Close()

    fmt.Printf("Loading %d sample blocks into %s...\n", numBlocks, dbPath)

    prevHash := "0"
    for i := 0; i < numBlocks; i++ {
        timestamp := time.Now().Unix() + int64(i*10)
        data := fmt.Sprintf("Transaction data for block %d", i)
        hash := blocks.ComputeHash(i, prevHash, data, timestamp)

        block := &blocks.Block{
            Height:    i,
            Hash:      hash,
            PrevHash:  prevHash,
            Data:      data,
            Timestamp: timestamp,
        }

        if err := storage.SaveBlock(block); err != nil {
            fmt.Printf("❌ Error saving block %d: %v\n", i, err)
            os.Exit(1)
        }

        fmt.Printf("✔ Block %d stored\n", i)
        prevHash = hash
    }

    fmt.Println("\n✅ Data loading complete!")
}

func viewBlock(dbPath, rpcURL string, jsonMode bool) {
    if flag.NArg() < 1 {
        fmt.Println("Usage: inspector -cmd block <height> [--rpc URL] [--json]")
        os.Exit(1)
    }
    
    height := 0
    fmt.Sscanf(flag.Arg(0), "%d", &height)
    
    var block *blocks.Block
    var err error
    
    if rpcURL != "" {
        fmt.Printf("Fetching block %d from RPC: %s\n", height, rpcURL)
        client := rpc.NewClient(rpcURL)
        block, err = client.FetchBlock(height)
        if err != nil {
            fmt.Printf("❌ Error fetching block via RPC: %v\n", err)
            os.Exit(1)
        }
    } else {
        storage, err := db.NewStorage(dbPath)
        if err != nil {
            fmt.Printf("❌ Error opening database: %v\n", err)
            os.Exit(1)
        }
        defer storage.Close()
        
        block, err = storage.LoadBlock(height)
        if err != nil {
            fmt.Printf("❌ Error loading block: %v\n", err)
            os.Exit(1)
        }
    }
    
    if jsonMode {
        data, _ := json.MarshalIndent(block, "", "  ")
        fmt.Println(string(data))
    } else {
        fmt.Printf("\n=== Block %d ===\n", block.Height)
        fmt.Printf("Hash:      %s\n", block.Hash)
        fmt.Printf("PrevHash:  %s\n", block.PrevHash)
        fmt.Printf("Timestamp: %s (Unix: %d)\n", time.Unix(block.Timestamp, 0).UTC(), block.Timestamp)
        fmt.Printf("Data:      %s\n\n", block.Data)
    }
}

func runScan(dbPath string, jsonMode bool) {
    storage, err := db.NewStorage(dbPath)
    if err != nil {
        fmt.Printf("❌ Error: %v\n", err)
        os.Exit(1)
    }
    defer storage.Close()

    result := errors.ScanErrors(storage, dbPath)
    errors.OutputScanResult(result, jsonMode)
}

func runCompare(db1Path, db2Path string, jsonMode bool) {
    storage1, err := db.NewStorage(db1Path)
    if err != nil {
        fmt.Printf("❌ Error opening Node1: %v\n", err)
        os.Exit(1)
    }
    defer storage1.Close()

    storage2, err := db.NewStorage(db2Path)
    if err != nil {
        fmt.Printf("❌ Error opening Node2: %v\n", err)
        os.Exit(1)
    }
    defer storage2.Close()

    result := errors.CompareNodes(storage1, storage2, db1Path, db2Path)
    errors.OutputComparisonResult(result, jsonMode)
}

func runConsensus(configPath string, jsonMode bool) {
    cfg, err := config.LoadConfig(configPath)
    if err != nil {
        fmt.Printf("❌ Error loading config: %v\n", err)
        os.Exit(1)
    }

    var nodes []consensus.NodeInfo
    for _, nodeConf := range cfg.Nodes {
        storage, err := db.NewStorage(nodeConf.DBPath)
        if err != nil {
            fmt.Printf("⚠️  Warning: Cannot open %s: %v\n", nodeConf.Name, err)
            continue
        }
        defer storage.Close()

        nodes = append(nodes, consensus.NodeInfo{
            Name:    nodeConf.Name,
            DBPath:  nodeConf.DBPath,
            Height:  storage.GetMaxHeight(),
            Storage: storage,
        })
    }

    if len(nodes) == 0 {
        fmt.Println("❌ Error: No valid nodes found")
        os.Exit(1)
    }

    result, err := consensus.AnalyzeConsensus(nodes)
    if err != nil {
        fmt.Printf("❌ Error analyzing consensus: %v\n", err)
        os.Exit(1)
    }

    consensus.OutputConsensusResult(result, jsonMode)
}

func runWatch(rpcURL string, interval int) {
    if rpcURL == "" {
        fmt.Println("❌ Error: --rpc flag is required for watch mode")
        fmt.Println("\nUsage: inspector -cmd watch --rpc http://localhost:8545 --interval 2")
        os.Exit(1)
    }
    
    watcher.Watch(rpcURL, interval)
}

func runFullReport(configPath, reportPath string) {
    fmt.Println("Generating comprehensive network report...")
    
    cfg, err := config.LoadConfig(configPath)
    if err != nil {
        fmt.Printf("❌ Error loading config: %v\n", err)
        os.Exit(1)
    }

    var nodes []consensus.NodeInfo
    for _, nodeConf := range cfg.Nodes {
        storage, err := db.NewStorage(nodeConf.DBPath)
        if err != nil {
            continue
        }
        defer storage.Close()

        nodes = append(nodes, consensus.NodeInfo{
            Name:    nodeConf.Name,
            DBPath:  nodeConf.DBPath,
            Height:  storage.GetMaxHeight(),
            Storage: storage,
        })
    }

    consensusResult, _ := consensus.AnalyzeConsensus(nodes)

    fullReport := &report.FullReport{
        Version:   version,
        Consensus: consensusResult,
    }

    if err := report.GenerateReport(reportPath, fullReport); err != nil {
        fmt.Printf("❌ Error generating report: %v\n", err)
        os.Exit(1)
    }
}

func printUsage() {
    fmt.Println("\n╔════════════════════════════════════════════════════════════════╗")
    fmt.Println("║          BHIV BLOCKCHAIN INSPECTOR CLI v1.0                   ║")
    fmt.Println("╚════════════════════════════════════════════════════════════════╝")
    fmt.Println("\n📋 COMMANDS:")
    fmt.Println("  load        Load sample blockchain data into LevelDB")
    fmt.Println("  block       View a specific block (LevelDB or RPC)")
    fmt.Println("  scan-errors Comprehensive blockchain error scanning")
    fmt.Println("  compare     Compare two blockchain databases")
    fmt.Println("  consensus   Multi-node consensus analysis (Day 2)")
    fmt.Println("  watch       Real-time node monitoring via RPC")
    fmt.Println("  report      Generate comprehensive JSON report (Day 3)")
    fmt.Println("  help        Show this help message")
    
    fmt.Println("\n💡 EXAMPLES:")
    
    fmt.Println("\n  Day 1 - Basic Operations:")
    fmt.Println("    inspector -cmd load -db ./data -blocks 50")
    fmt.Println("    inspector -cmd block 10 --db ./data")
    fmt.Println("    inspector -cmd block 10 --rpc http://localhost:8545")
    fmt.Println("    inspector -cmd scan-errors --db ./data")
    fmt.Println("    inspector -cmd compare -db1 ./node1 -db2 ./node2")
    
    fmt.Println("\n  Day 2 - Consensus Analysis:")
    fmt.Println("    inspector -cmd consensus --config nodes.json")
    fmt.Println("    inspector -cmd consensus --config nodes.json --json")
    
    fmt.Println("\n  Day 3 - Advanced:")
    fmt.Println("    inspector -cmd watch --rpc http://localhost:8545 --interval 2")
    fmt.Println("    inspector -cmd report --config nodes.json --report network-report.json")
    
    fmt.Println("\n🔧 OPTIONS:")
    fmt.Println("  --db         Path to LevelDB database")
    fmt.Println("  --rpc        RPC endpoint URL")
    fmt.Println("  --config     Network config file (JSON)")
    fmt.Println("  --json       Output in JSON format")
    fmt.Println("  --interval   Watch mode polling interval")
    fmt.Println("  --report     Report output path")
    fmt.Println()
}

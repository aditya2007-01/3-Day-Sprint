package main

import (
    "encoding/json"
    "flag"
    "fmt"
    "io/ioutil"
    "log"
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

// DAY 1: JSON OUTPUT STRUCT [file:15]
type BlockJSON struct {
    Height       int    `json:"height"`
    BlockHash    string `json:"blockHash"`
    PreviousHash string `json:"previousHash"`
    TxCount      int    `json:"txCount"`
}

var verboseFlag bool

func main() {
    // DAY 1 PDF-REQUIRED FLAGS [file:15]
    path := flag.String("path", "", "leveldb-path (Day 1)")
    compare1 := flag.String("compare1", "", "path1 for comparison")
    compare2 := flag.String("compare2", "", "path2 for comparison")
    height := flag.Int("height", 0, "block height n")
    jsonOutput := flag.Bool("json", false, "export structured JSON output")
    verbose := flag.Bool("verbose", false, "verbose logging")
    quiet := flag.Bool("quiet", false, "quiet mode")
    
    // EXISTING FLAGS
    dbPath := flag.String("db", "./leveldb-data", "Path to LevelDB database")
    db1Path := flag.String("db1", "./node1-data", "Path to first database")
    db2Path := flag.String("db2", "./node2-data", "Path to second database")
    cmd := flag.String("cmd", "help", "Command: load, block, scan-errors, compare, consensus, watch, report")
    numBlocks := flag.Int("blocks", 10, "Number of blocks to load")
    showVersion := flag.Bool("version", false, "Show version")
    rpcURL := flag.String("rpc", "", "RPC endpoint URL")
    watchInterval := flag.Int("interval", 2, "Watch mode interval in seconds")
    configPath := flag.String("config", "nodes.json", "Path to network config file")
    reportPath := flag.String("report", "inspector-report.json", "Output path for report")
    
    flag.Parse()

    // DAY 1: LOGGING CONTROL [file:15]
    verboseFlag = *verbose
    if *verbose {
        log.SetFlags(log.LstdFlags | log.Lshortfile)
        log.Println("✓ Verbose mode enabled")
    } else if *quiet {
        log.SetOutput(ioutil.Discard)
    }

    if *showVersion {
        fmt.Printf("BHIV Chain Inspector v%s\n", version)
        return
    }

    // DAY 1: SHORTCUT ROUTES [file:15]
    if *path != "" && *height > 0 {
        viewBlockDay1(*path, *rpcURL, *height, *jsonOutput)
        return
    }
    
    if *compare1 != "" && *compare2 != "" {
        compareNodesDay1(*compare1, *compare2, *jsonOutput)
        return
    }

    // EXISTING COMMANDS
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

// DAY 1: NEW FUNCTION [file:15]
func viewBlockDay1(dbPath, rpcURL string, height int, jsonMode bool) {
    var block *blocks.Block
    var err error
    
    if rpcURL != "" {
        if verboseFlag {
            log.Printf("Fetching block %d from RPC: %s", height, rpcURL)
        }
        client := rpc.NewClient(rpcURL)
        block, err = client.FetchBlock(height)
        if err != nil {
            errors.FormatError("BLOCK_FETCH_FAILED", err.Error(), height)
            return
        }
    } else {
        if verboseFlag {
            log.Printf("Opening database: %s", dbPath)
        }
        storage, err := db.NewStorage(dbPath)
        if err != nil {
            errors.FormatError("DB_OPEN_FAILED", err.Error(), height)
            return
        }
        defer storage.Close()
        
        if verboseFlag {
            log.Printf("Loading block at height: %d", height)
        }
        block, err = storage.LoadBlock(height)
        if err != nil {
            errors.FormatError("BLOCK_FETCH_FAILED", err.Error(), height)
            return
        }
    }

    output := BlockJSON{
        Height:       block.Height,
        BlockHash:    block.Hash,
        PreviousHash: block.PrevHash,
        TxCount:      1,
    }

    if jsonMode {
        encoder := json.NewEncoder(os.Stdout)
        encoder.SetIndent("", "  ")
        encoder.Encode(output)
    } else {
        fmt.Printf("\n=== Block %d ===\n", output.Height)
        fmt.Printf("Hash:      %s\n", output.BlockHash)
        fmt.Printf("PrevHash:  %s\n", output.PreviousHash)
        fmt.Printf("TxCount:   %d\n\n", output.TxCount)
    }
}

// DAY 1: NEW FUNCTION [file:15]
func compareNodesDay1(path1, path2 string, jsonMode bool) {
    if verboseFlag {
        log.Printf("Opening node 1: %s", path1)
    }
    storage1, err := db.NewStorage(path1)
    if err != nil {
        errors.FormatError("DB_OPEN_FAILED", fmt.Sprintf("Node1: %v", err), 0)
        return
    }
    defer storage1.Close()

    if verboseFlag {
        log.Printf("Opening node 2: %s", path2)
    }
    storage2, err := db.NewStorage(path2)
    if err != nil {
        errors.FormatError("DB_OPEN_FAILED", fmt.Sprintf("Node2: %v", err), 0)
        return
    }
    defer storage2.Close()

    if verboseFlag {
        log.Println("Starting node comparison...")
    }
    result := errors.CompareNodes(storage1, storage2, path1, path2)
    errors.OutputComparisonResult(result, jsonMode)
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
    fmt.Println("║          BHIV BLOCKCHAIN INSPECTOR CLI v1.0 - DAY 1           ║")
    fmt.Println("╚════════════════════════════════════════════════════════════════╝")
    
    fmt.Println("\n📋 DAY 1 COMMANDS (PDF Required):")
    fmt.Println("  inspector --path <db> --height <n> [--json] [--verbose]")
    fmt.Println("  inspector --compare1 <path1> --compare2 <path2> [--json]")
    
    fmt.Println("\n💡 DAY 1 EXAMPLES (Record These):")
    fmt.Println("  inspector --path ./data --height 10 --json")
    fmt.Println("  inspector --path ./data --height 5 --verbose")
    fmt.Println("  inspector --compare1 ./node1 --compare2 ./node2 --json")
    
    fmt.Println("\n📋 ORIGINAL COMMANDS:")
    fmt.Println("  load        Load sample blocks")
    fmt.Println("  block       View specific block")
    fmt.Println("  scan-errors Scan for errors")
    fmt.Println("  compare     Compare two nodes")
    fmt.Println("  consensus   Consensus analysis")
    fmt.Println("  watch       Real-time monitoring")
    fmt.Println("  report      Generate report")
    
    fmt.Println("\n🔧 FLAGS:")
    fmt.Println("  --path       leveldb-path")
    fmt.Println("  --height     block height")
    fmt.Println("  --compare1   first node path")
    fmt.Println("  --compare2   second node path")
    fmt.Println("  --json       JSON output")
    fmt.Println("  --verbose    verbose mode")
    fmt.Println("  --quiet      quiet mode")
    fmt.Println()
}

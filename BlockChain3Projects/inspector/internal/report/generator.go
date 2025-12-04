package report

import (
    "encoding/json"
    "fmt"
    "os"
    "time"

    "inspector/internal/consensus"
    "inspector/internal/errors"
)

type FullReport struct {
    Timestamp       string                      `json:"timestamp"`
    Version         string                      `json:"version"`
    ErrorScan       *errors.ErrorScanResult     `json:"error_scan,omitempty"`
    Comparison      *errors.ComparisonResult    `json:"comparison,omitempty"`
    Consensus       *consensus.ConsensusResult  `json:"consensus,omitempty"`
}

func GenerateReport(outputPath string, report *FullReport) error {
    report.Timestamp = time.Now().Format("2006-01-02 15:04:05")
    
    jsonData, err := json.MarshalIndent(report, "", "  ")
    if err != nil {
        return fmt.Errorf("failed to marshal report: %w", err)
    }

    if err := os.WriteFile(outputPath, jsonData, 0644); err != nil {
        return fmt.Errorf("failed to write report: %w", err)
    }

    fmt.Printf("✅ Report saved to: %s\n", outputPath)
    return nil
}

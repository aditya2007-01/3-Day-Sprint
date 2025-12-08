package blocks

import (
    "testing"
)

// TEST 1: Block deserialization [file:15]
func TestComputeHash(t *testing.T) {
    hash := ComputeHash(1, "prev", "data", 1234567890)
    
    if hash == "" {
        t.Error("Hash should not be empty")
    }
    
    if len(hash) != 64 { // SHA256 produces 64 hex chars
        t.Errorf("Expected hash length 64, got %d", len(hash))
    }
}

// TEST 2: Hash verification [file:15]
func TestHashVerification(t *testing.T) {
    prevHash := "0"
    data := "test data"
    timestamp := int64(1234567890)
    
    hash1 := ComputeHash(1, prevHash, data, timestamp)
    hash2 := ComputeHash(1, prevHash, data, timestamp)
    
    if hash1 != hash2 {
        t.Error("Same inputs should produce same hash")
    }
    
    hash3 := ComputeHash(2, prevHash, data, timestamp)
    if hash1 == hash3 {
        t.Error("Different heights should produce different hashes")
    }
}

func TestBlockStructure(t *testing.T) {
    block := Block{
        Height:    10,
        Hash:      "abc123",
        PrevHash:  "def456",
        Data:      "transaction data",
        Timestamp: 1234567890,
    }
    
    if block.Height != 10 {
        t.Errorf("Expected height 10, got %d", block.Height)
    }
    
    if block.Hash == "" {
        t.Error("Hash should not be empty")
    }
}

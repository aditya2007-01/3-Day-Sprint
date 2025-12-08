package blocks

import (
    "testing"
)

// TEST 5: Merkle root computation (simplified) [file:15]
func TestHashConsistency(t *testing.T) {
    data1 := "block data 1"
    data2 := "block data 2"
    
    hash1 := ComputeHash(1, "0", data1, 1000)
    hash2 := ComputeHash(1, "0", data1, 1000)
    hash3 := ComputeHash(1, "0", data2, 1000)
    
    if hash1 != hash2 {
        t.Error("Identical inputs should produce identical hashes")
    }
    
    if hash1 == hash3 {
        t.Error("Different data should produce different hashes")
    }
}

func TestHashDeterminism(t *testing.T) {
    tests := []struct {
        height    int
        prevHash  string
        data      string
        timestamp int64
    }{
        {1, "0", "genesis", 1000},
        {2, "abc", "tx1", 2000},
        {3, "def", "tx2", 3000},
    }
    
    for _, tt := range tests {
        hash1 := ComputeHash(tt.height, tt.prevHash, tt.data, tt.timestamp)
        hash2 := ComputeHash(tt.height, tt.prevHash, tt.data, tt.timestamp)
        
        if hash1 != hash2 {
            t.Errorf("Hash not deterministic for height %d", tt.height)
        }
    }
}

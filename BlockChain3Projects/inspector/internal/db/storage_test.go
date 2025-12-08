package db

import (
    "os"
    "testing"
    
    "inspector/internal/blocks"
)

// TEST 3: Database operations [file:15]
func TestStorageOpenClose(t *testing.T) {
    testPath := "./test_db"
    defer os.RemoveAll(testPath)
    
    storage, err := NewStorage(testPath)
    if err != nil {
        t.Fatalf("Failed to create storage: %v", err)
    }
    
    if storage == nil {
        t.Fatal("Storage should not be nil")
    }
    
    storage.Close()
}

func TestSaveAndLoadBlock(t *testing.T) {
    testPath := "./test_db_save"
    defer os.RemoveAll(testPath)
    
    storage, err := NewStorage(testPath)
    if err != nil {
        t.Fatalf("Failed to create storage: %v", err)
    }
    defer storage.Close()
    
    block := &blocks.Block{
        Height:    5,
        Hash:      "testhash123",
        PrevHash:  "prevhash456",
        Data:      "test data",
        Timestamp: 1234567890,
    }
    
    err = storage.SaveBlock(block)
    if err != nil {
        t.Fatalf("Failed to save block: %v", err)
    }
    
    loadedBlock, err := storage.LoadBlock(5)
    if err != nil {
        t.Fatalf("Failed to load block: %v", err)
    }
    
    if loadedBlock.Height != block.Height {
        t.Errorf("Expected height %d, got %d", block.Height, loadedBlock.Height)
    }
    
    if loadedBlock.Hash != block.Hash {
        t.Errorf("Expected hash %s, got %s", block.Hash, loadedBlock.Hash)
    }
}

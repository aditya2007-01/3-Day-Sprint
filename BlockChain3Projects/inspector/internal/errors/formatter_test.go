package errors

import (
    "testing"
)

// TEST 4: Error formatter module [file:15]
func TestBlockErrorStructure(t *testing.T) {
    err := BlockError{
        Height: 10,
        Code:   "TEST_ERROR",
        Msg:    "Test message",
    }
    
    if err.Height != 10 {
        t.Errorf("Expected height 10, got %d", err.Height)
    }
    
    if err.Code != "TEST_ERROR" {
        t.Errorf("Expected code TEST_ERROR, got %s", err.Code)
    }
    
    if err.Msg == "" {
        t.Error("Message should not be empty")
    }
}

func TestErrorCodes(t *testing.T) {
    testCases := []struct {
        code   string
        msg    string
        height int
    }{
        {"DB_OPEN_FAILED", "Database not found", 0},
        {"BLOCK_FETCH_FAILED", "Block missing", 10},
        {"HASH_MISMATCH", "Invalid hash", 5},
    }
    
    for _, tc := range testCases {
        err := BlockError{
            Height: tc.height,
            Code:   tc.code,
            Msg:    tc.msg,
        }
        
        if err.Code != tc.code {
            t.Errorf("Expected code %s, got %s", tc.code, err.Code)
        }
    }
}

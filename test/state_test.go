package test

import (
	"os"
	"testing"
)
import "demoBlockChain/database"

func setEnv() bool {
	err := os.Setenv("DATABASE_PATH", "./_testdata")
	if err != nil {
		return false
	}
	return true
}

func init() {
	if !setEnv() {
		panic("set env failed")
	}
}

func TestNewStateFromDisk(t *testing.T) {
	state, err := database.NewStateFromDisk()
	if err != nil {
		t.Errorf("NewStateFromDisk() error: %v", err)
	}
	if len(state.Balances) == 0 {
		t.Errorf("NewStateFromDisk() error: %v", "state.Balances is empty")
	}
}

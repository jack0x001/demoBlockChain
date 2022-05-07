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

func TestAddTx(t *testing.T) {
	state, _ := database.NewState("")

	state.Balances["zhou"] = 100
	state.Balances["li"] = 200

	t.Run("normal case", func(t *testing.T) {
		t.Helper()

		tx := database.Tx{
			From:  "zhou",
			To:    "li",
			Value: 5,
			Data:  "",
		}
		err := state.AddTx(tx)
		if err != nil {
			t.Errorf("add failed: %s", err)
		}

		if state.Balances["zhou"] != 95 {
			t.Errorf("add failed: the balance should be %d", state.Balances["zhou"])
		}
		if state.Balances["li"] != 205 {
			t.Errorf("add failed: the balance should be %d", state.Balances["li"])
		}
	})

	t.Run("invalid case", func(t *testing.T) {
		t.Helper()
		tx := database.Tx{
			From:  "zhou",
			To:    "li",
			Value: 500,
			Data:  "",
		}
		err := state.AddTx(tx)
		if err == nil {
			t.Errorf("invalid case failed, the tx should be failed")
		}
	})

}

func TestPersist(t *testing.T) {
	state, err := database.NewState("./_testdata/tx_test.db")
	defer state.CloseDB()

	state.Balances["zhouyh"] = 100
	state.Balances["li"] = 200

	_ = state.AddTx(database.Tx{
		From:  "zhouyh",
		To:    "li",
		Value: 5,
		Data:  "",
	})

	_ = state.AddTx(database.Tx{
		From:  "li",
		To:    "zhouyh",
		Value: 20,
		Data:  "",
	})

	if len(state.TxMemPool) != 2 {
		t.Errorf("add failed: the tx pool length should be 2, but got %d", len(state.TxMemPool))
	}

	err = state.Persist()
	if err != nil {
		t.Errorf("persist failed, %s", err.Error())
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
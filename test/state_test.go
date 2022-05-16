package test

import (
	"demoBlockChain/database"
	"testing"
)

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

func TestPersistAndLoad(t *testing.T) {
	state, err := database.NewStateFromDisk("./_testdata")
	if err != nil {
		t.Errorf("new state failed: %s", err)
	}
	defer state.CloseDB()

	oldHeight := state.LatestBlock.Header.Height
	tx := database.Tx{
		From:  "zhou",
		To:    "li",
		Value: 500,
		Data:  "",
	}
	err = state.AddTx(tx)
	if err == nil {
		t.Errorf("invalid case failed, the tx should be failed")
	}

	_, err = state.Persist()
	if err != nil {
		t.Errorf("persist failed, %s", err.Error())
	}

	state, err = database.NewStateFromDisk("./_testdata")
	if err != nil {
		t.Errorf("new state failed: %s", err)
	}

	if state.LatestBlock.Header.Height != oldHeight+1 {
		t.Errorf("block height, the height should be %d + 1 ", oldHeight)
	}
}

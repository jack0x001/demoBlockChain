package database

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// State 表示状态机的当前状态, 用户余额以及交易列表
type State struct {
	Balances  map[Account]uint
	txMemPool []Tx

	dbFile *os.File
}

func (s State) apply(tx Tx) error {
	if tx.IsReward() {
		s.Balances[tx.To] += tx.Value
		return nil
	}

	if s.Balances[tx.From] < tx.Value {
		return fmt.Errorf("余额不足")
	}

	s.Balances[tx.From] -= tx.Value
	s.Balances[tx.To] += tx.Value
	return nil
}

// NewStateFromDisk 从磁盘读取并生成新的状态
func NewStateFromDisk() (state *State, err error) {

	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	//读取创世文件
	genesisFilePath := filepath.Join(cwd, "database", "genesis.json")
	genesis, err := loadGenesis(genesisFilePath)
	if err != nil {
		return nil, err
	}
	balances := make(map[Account]uint)
	for account, value := range genesis.Balances {
		balances[account] = value
	}

	//读取交易记录
	txFilePath := filepath.Join(cwd, "database", "tx.db")
	txFile, err := os.OpenFile(txFilePath, os.O_APPEND|os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(txFile)
	state = &State{balances, make([]Tx, 0), txFile}
	for scanner.Scan() {
		if err = scanner.Err(); err != nil {
			return nil, err
		}
		var tx Tx
		err = json.Unmarshal(scanner.Bytes(), &tx)
		if err != nil {
			return nil, err
		}
		err = state.apply(tx)
		if err != nil {
			return nil, err
		}
	}

	return state, nil
}

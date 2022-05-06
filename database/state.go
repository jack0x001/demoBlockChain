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

// 验证交易, 如果验证通过, 则将使交易生效
func (s State) apply(tx Tx) error {
	if tx.IsReward() {
		s.Balances[tx.To] += tx.Value
		return nil
	}

	if s.Balances[tx.From] < tx.Value {
		return fmt.Errorf("余额不足, 当前余额 %d 小于所需值 %d", s.Balances[tx.From], tx.Value)
	}

	s.Balances[tx.From] -= tx.Value
	s.Balances[tx.To] += tx.Value
	return nil
}

// NewStateFromDisk 从磁盘读取并生成新的状态, 用于恢复状态机
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
	//defer func(txFile *os.File) {
	//	err := txFile.Close()
	//	if err != nil {
	//		log.Println(err)
	//	}
	//}(txFile)
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

// Add : 添加一条交易到MemPool
func (s State) Add(tx Tx) error {
	if err := s.apply(tx); err != nil {
		return err
	}
	s.txMemPool = append(s.txMemPool, tx)
	return nil
}

// Persist : 持久化状态到磁盘
func (s State) Persist() error {

	//logic:
	//目的: 将s.txMemPool写到磁盘
	//但是, 写磁盘的时候s.txMemPool有可能正在被Add进行append
	//所以:
	//1 ,  使用s.txMemPool的副本进行磁盘写入
	//2 ,  写入成功, 则将其从s.txMemPool中删除

	/*

		          s.txMemPool

		         ┌─┐┌─┐┌─┐┌─┐┌─┐┌─┐┌─┐┌─┐  ◀─────adding (append)──────
		    ─ ─ ▶└─┘└─┘└─┘└─┘└─┘└─┘└─┘└─┘
		   │
		             copied
		   │     ┌─┐┌─┐┌─┐┌─┐┌─┐
		         └─┘└─┘└─┘└─┘└─┘
		   │      │
		delete    │
		   │      └────┐
		               │
		   │           ▼
		      ┌────────────────┐
		   └ ─│   write disk   │
		      └────────────────┘

	*/

	memoryPool := make([]Tx, len(s.txMemPool))
	copy(memoryPool, s.txMemPool)

	for i := 0; i < len(memoryPool); i++ {
		jsonBytes, err := json.Marshal(memoryPool[i])
		if err != nil {
			return err
		}
		jsonBytesWithNewLine := append(jsonBytes, '\n')

		_, err = s.dbFile.Write(jsonBytesWithNewLine)
		if err != nil {
			return err
		}

		//remove s.txMemPool[0]
		s.txMemPool = s.txMemPool[1:]
	}
	return nil
}

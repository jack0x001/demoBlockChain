package database

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// State 表示状态机的当前状态, 用户余额以及交易列表
type State struct {
	Balances  map[Account]uint
	TxMemPool []Tx

	dbFile *os.File
}

func NewState(dbFilePath string) (*State, error) {
	s := &State{}
	s.Balances = make(map[Account]uint)
	s.TxMemPool = make([]Tx, 0)

	//if dbFilePath is not empty, then open the file
	if dbFilePath != "" {
		dbFile, err := os.OpenFile(dbFilePath, os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			return nil, err
		}
		s.dbFile = dbFile
	} else {
		s.dbFile = nil
	}

	return s, nil
}

// 验证交易, 如果验证通过, 则将使交易生效
func (s *State) apply(tx Tx) error {
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

	// get environment variable of database path
	dbPath := os.Getenv("DATABASE_PATH")
	//读取创世文件
	genesisFilePath := filepath.Join(dbPath, "genesis.json")
	genesis, err := loadGenesis(genesisFilePath)
	if err != nil {
		return nil, err
	}
	balances := make(map[Account]uint)
	for account, value := range genesis.Balances {
		balances[account] = value
	}

	//读取交易记录
	txFilePath := filepath.Join(dbPath, "tx.db")
	txFile, err := os.OpenFile(txFilePath, os.O_APPEND|os.O_RDWR, 0600)
	//defer func(txFile *os.File) {
	//	err := txFile.CloseDB()
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

// AddTx : 添加一条交易
func (s *State) AddTx(tx Tx) error {
	if err := s.apply(tx); err != nil {
		return err
	}
	s.TxMemPool = append(s.TxMemPool, tx)
	return nil
}

// Persist : 持久化状态到磁盘
func (s *State) Persist() error {

	//logic:
	//目的: 将s.txMemPool写到磁盘
	//但是, 写磁盘的时候s.txMemPool有可能正在被Add进行append
	//所以:
	//1 ,  使用s.txMemPool的副本进行磁盘写入
	//2 ,  写入成功, 则将其从s.txMemPool中删除

	/*

		          s.TxMemPool

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

	memoryPool := make([]Tx, len(s.TxMemPool))
	copy(memoryPool, s.TxMemPool)

	for i := 0; i < len(memoryPool); i++ {
		jsonBytes, err := json.Marshal(memoryPool[i])
		if err != nil {
			return err
		}
		jsonBytesWithNewLine := append(jsonBytes, '\n')

		if s.dbFile == nil {
			return errors.New("dbFile of state is nil")
		}
		_, err = s.dbFile.Write(jsonBytesWithNewLine)
		if err != nil {
			return err
		}

		//remove s.TxMemPool[0]
		s.TxMemPool = s.TxMemPool[1:]
	}
	return nil
}

func (s *State) CloseDB() {
	err := s.dbFile.Close()
	if err != nil {
		return
	}
}

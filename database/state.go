package database

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// State 表示状态机的当前状态, 用户余额以及交易列表
type State struct {
	Balances  map[Account]uint
	TxMemPool []Tx // 交易池

	dbFile        *os.File
	lastBlockHash HashCode // 最后一个区块的哈希值
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
		return fmt.Errorf("%q 余额不足, 当前余额 %d 小于所需值 %d", tx.From, s.Balances[tx.From], tx.Value)
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

	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(txFile)
	state = &State{balances, make([]Tx, 0), txFile, HashCode{}}

	for scanner.Scan() {
		if err = scanner.Err(); err != nil {
			return nil, err
		}

		var blockFS BlockFS
		err = json.Unmarshal(scanner.Bytes(), &blockFS)
		if err != nil {
			return nil, err
		}
		err = state.applyBlock(blockFS.Value)
		if err != nil {
			return nil, err
		}

		state.lastBlockHash = blockFS.Key
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

func (s *State) applyBlock(block Block) error {
	for _, tx := range block.Transactions {
		if err := s.apply(tx); err != nil {
			return err
		}
	}
	return nil
}

// Persist : 持久化状态到磁盘
func (s *State) Persist() (HashCode, error) {

	block := NewBlock(s.lastBlockHash, uint64(time.Now().Unix()), s.TxMemPool)
	blockHash, err := block.Hash()
	if err != nil {
		return HashCode{}, err
	}

	blockFS := BlockFS{
		blockHash,
		*block,
	}

	blockFSJson, err := json.Marshal(blockFS)
	if err != nil {
		return HashCode{}, err
	}

	_, err = s.dbFile.Write(append(blockFSJson, '\n'))
	if err != nil {
		return HashCode{}, err
	}
	s.lastBlockHash = blockHash
	s.TxMemPool = make([]Tx, 0)

	return blockHash, nil

}

func (s *State) CloseDB() {
	err := s.dbFile.Close()
	if err != nil {
		return
	}
}

package database

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
)

type HashCode [32]byte

// MarshalText 重写序列化方法
func (h HashCode) MarshalText() ([]byte, error) {
	return []byte(hex.EncodeToString(h[:])), nil
}

// UnmarshalText 重写反序列化方法
func (h *HashCode) UnmarshalText(data []byte) error {
	_, err := hex.Decode(h[:], data)
	return err
}

type BlockHeader struct {
	Parent HashCode `json:"parent"`
	Time   uint64   `json:"time"`
}

type Block struct {
	Header       BlockHeader `json:"header"`
	Transactions []Tx        `json:"payload"`
}

// Hash returns the HashCode of the Block
func (b *Block) Hash() (HashCode, error) {

	j, err := json.Marshal(b)
	if err != nil {
		return HashCode{}, err
	}
	return sha256.Sum256(j), nil
}

func NewBlock(parent HashCode, time uint64, transactions []Tx) *Block {
	return &Block{
		Header: BlockHeader{
			Parent: parent,
			Time:   time,
		},
		Transactions: transactions,
	}
}

type BlockFS struct {
	Key   HashCode `json:"hash"`
	Value Block    `json:"block"`
}

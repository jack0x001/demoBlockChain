package database

import (
	"encoding/json"
	"io/ioutil"
)

type genesis struct {
	Balances map[Account]uint `json:"balances"`
}

// 读取并解析genesis.json文件
func loadGenesis(path string) (genesis, error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return genesis{}, err
	}
	var g genesis
	err = json.Unmarshal(bytes, &g)
	if err != nil {
		return genesis{}, err
	}
	return g, nil
}

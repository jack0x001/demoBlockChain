package database

import (
	"encoding/json"
	"io/ioutil"
	"os"
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

func writeGenesis(path string) error {
	genesisJson := `{
  "genesis_time": "2022-05-12T00:00:00.000000000Z",
  "chain_id": "the-demo-blockchain-ledger",
  "balances": {
    "zhouyh": 1000000
  }
}`

	return os.WriteFile(path, []byte(genesisJson), 0644)
}

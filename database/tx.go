package database

type Account string

func NewAccount(name string) Account {
	return Account(name)
}

// Tx 表示一次交易.
type Tx struct {
	From  Account `json:"from"`
	To    Account `json:"to"`
	Value uint    `json:"value"`
	Data  string  `json:"data"`
}

// NewTx 新生成一个交易.
func NewTx(from Account, to Account, value uint, data string) Tx {
	return Tx{
		From:  from,
		To:    to,
		Value: value,
		Data:  data,
	}
}

func (t Tx) IsReward() bool {
	return t.Data == "reward"
}

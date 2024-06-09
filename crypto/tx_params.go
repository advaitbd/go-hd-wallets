package crypto

type DataType int

const (
	Transaction DataType = iota + 1
	TypedData
	PersonalMessage
	TypedTransaction
)

type TxParams struct {
	Nonce    uint64 `json:"nonce"`
	To       string `json:"to"`
	From     string `json:"from"`
	Data     string `json:"data"`
	GasLimit uint64 `json:"gasLimit,string"`
	GasPrice uint64 `json:"gasPrice,string"`
	Value    uint64 `json:"value"`
	ChainID  uint64 `json:"chainId"`
}

package storage

var Db Storage // polymorphic

type Storage interface {
	GetLastProcessedBlock() (uint64, error)
	SetLastProcessedBlock(value uint64) error

	InsertWallet(wallet Wallet) error
	GetWallets() ([]Wallet, error)

	InsertTransaction(address string, transaction Transaction) error
	GetTransactions(address string) ([]Transaction, error)
}

type Wallet struct {
	Address string
}

type Transaction struct {
	FromAddress string
	ToAddress   string
	Value       uint64
	BlockHeight uint64
}

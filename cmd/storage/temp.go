package storage

import (
	"errors"
	"sync"
)

type TempStorage struct {
	hashmap sync.Map // thread-safe
}

func NewTempStorage() *TempStorage {
	return &TempStorage{}
}

func (s *TempStorage) GetLastProcessedBlock() (uint64, error) {
	value, ok := s.hashmap.Load("lastProcessedBlock")

	if !ok {
		return 0, nil
	}

	if valueInt, ok := value.(uint64); ok {
		return valueInt, nil
	}

	return 0, errors.New("invalid lastProcessdBlock value")
}

func (s *TempStorage) SetLastProcessedBlock(value uint64) error {
	s.hashmap.Store("lastProcessedBlock", value)
	return nil
}

func (s *TempStorage) InsertWallet(wallet Wallet) error {
	wallets, _ := s.hashmap.LoadOrStore("wallets", []Wallet{})
	walletsSlice := wallets.([]Wallet)
	s.hashmap.Store("wallets", append(walletsSlice, wallet))
	return nil
}

func (s *TempStorage) GetWallets() ([]Wallet, error) {
	wallets, ok := s.hashmap.Load("wallets")
	if !ok {
		return []Wallet{}, nil
	}
	return wallets.([]Wallet), nil
}

func (s *TempStorage) InsertTransaction(address string, transaction Transaction) error {
	transactions, _ := s.hashmap.LoadOrStore(address, []Transaction{})
	transactionsSlice := transactions.([]Transaction)
	s.hashmap.Store(address, append(transactionsSlice, transaction))
	return nil
}

func (s *TempStorage) GetTransactionsByAddress(address string) ([]Transaction, error) {
	transactions, ok := s.hashmap.Load(address)
	if !ok {
		return []Transaction{}, nil
	}
	return transactions.([]Transaction), nil
}

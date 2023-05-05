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

func (s *TempStorage) GetTransactions(address string) ([]Transaction, error) {
	return nil, nil
}

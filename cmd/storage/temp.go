package storage

import (
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

	return 0, nil
}

func (s *TempStorage) SetLastProcessedBlock(value uint64) error {
	s.hashmap.Store("lastProcessedBlock", value)
	return nil
}

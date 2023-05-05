package storage

type RedisStorage struct {
	// go-redis (thread-safe)
}

func NewRedisStorage() *RedisStorage {
	panic("txmond: redis storage backend is not implemented")
}

func (s *RedisStorage) GetLastProcessedBlock() (uint64, error) {
	panic("txmond: redis storage backend is not implemented")
}

func (s *RedisStorage) SetLastProcessedBlock(value uint64) error {
	panic("txmond: redis storage backend is not implemented")
}

func (s *RedisStorage) InsertWallet(wallet Wallet) error {
	panic("txmond: redis storage backend is not implemented")
}

func (s *RedisStorage) GetWallets() ([]Wallet, error) {
	panic("txmond: redis storage backend is not implemented")
}

func (s *RedisStorage) InsertTransaction(adress string, transaction Transaction) error {
	panic("txmond: redis storage backend is not implemented")
}

func (s *RedisStorage) GetTransactionsByAddress(address string) ([]Transaction, error) {
	panic("txmond: redis storage backend is not implemented")
}

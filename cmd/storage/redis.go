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

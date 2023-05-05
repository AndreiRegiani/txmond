package storage

var Db Storage // polymorphic

type Storage interface {
	GetLastProcessedBlock() (uint64, error)
	SetLastProcessedBlock(value uint64) error
}

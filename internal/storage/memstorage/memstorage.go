package memstorage

import "sync"

type MemStorage struct {
	metrics sync.Map
}

func New() (s *MemStorage) {
	return &MemStorage{
		metrics: sync.Map{},
	}
}

func (s *MemStorage) Set(metric string, value float64) {
	s.metrics.Store(metric, value)
}

func (s *MemStorage) Update(metric string, value int64) {
	s.metrics.LoadOrStore(metric, value)
}

package storage

type Storage interface {
	Set(name string, value float64)
	Update(name string, value int64)
}

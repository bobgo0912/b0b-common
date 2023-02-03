package sql

type Pagination[T any] struct {
	Page  uint64
	Size  uint64
	Total uint64
	Data  []*T
}

package cache

type Cacher interface {
	Set([]byte, []byte, int64) error
	Has([]byte) bool
	Get([]byte) ([]byte, error)
	Delete([]byte) error
}

// Provides inmemory storage to keep audiohash: text pairs
package memstore

type MemStore struct {
	Tuple map[string]string
}

// Save tuple to RAM
func Save(hash string, text string) (string, error) {

	return "", nil
}

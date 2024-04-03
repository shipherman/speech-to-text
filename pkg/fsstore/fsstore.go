// Package fsstore provides functionality
// to save files to local store
//

package fsstore

import "os"

// FSStore implements store interface in models
type FSStore struct {
	Path string
}

// NewFSStore creates an instance of FSStore
func NewFSStore(path string) *FSStore {
	return &FSStore{Path: path}
}

// Configure changes path to local store
func (f *FSStore) Configure(path string) error {
	f.Path = path
	return nil
}

// Save func saves data to a local path
func (f *FSStore) Save(name string, data []byte) error {
	filePath := f.Path + "/" + name
	osFile, err := os.Create(filePath)
	if err != nil {
		return err
	}
	_, err = osFile.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func (f *FSStore) Get(name string) ([]byte, error) {
	b, err := os.ReadFile(name)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (f *FSStore) GetStorePath() string {
	return f.Path
}

func (f *FSStore) Close() error {
	return nil
}

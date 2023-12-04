// Package fsstore provides functionality
// to save files to local store
//

package fsstore

import "os"

// FSStore implements store interface in models
type FSStore struct {
	Path string
}

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

func (f *FSStore) Delete(name string) error {
	return nil
}

package fsstore

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const path = "/tmp/store"

func TestMain(m *testing.M) {
	os.Mkdir(path, os.ModePerm)

	code := m.Run()

	os.RemoveAll(path)

	os.Exit(code)
}

func TestNewFSStore(t *testing.T) {

	s := NewFSStore(path)
	assert.Equal(t, path, s.Path)

	err := s.Save("file", []byte("file"))
	assert.NoError(t, err)

	b, err := s.Get(path + "/" + "file")
	assert.Equal(t, []byte("file"), b)
	assert.NoError(t, err)
}

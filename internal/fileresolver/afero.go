package fileresolver

import "github.com/spf13/afero"

// AferoFile wraps an afero file into an ABI.File
type AferoFile struct {
	afero.File
}

// Flush runs Sync.
func (a AferoFile) Flush() error { return a.Sync() }

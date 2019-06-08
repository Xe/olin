package fileresolver

import "within.website/olin/internal/abi"

// Zero is a file that does nothing.
//
// For more information, please see the spec here: https://github.com/CommonWA/cwa-spec/blob/master/schemes/zero.md
func Zero() abi.File {
	return zeroFile{}
}

type zeroFile struct{}

func (zeroFile) Write(p []byte) (int, error) { return len(p), nil }
func (zeroFile) Flush() error                { return nil }
func (zeroFile) Close() error                { return nil }
func (zeroFile) Name() string                { return "zero" }

func (zeroFile) Read(p []byte) (int, error) {
	for i := range p {
		p[i] = 0
	}

	return len(p), nil
}

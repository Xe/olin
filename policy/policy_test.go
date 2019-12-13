package policy

import (
	"io/ioutil"
	"path/filepath"
	"testing"
)

func TestParse(t *testing.T) {
	for _, file := range []string{"gitea.policy"} {
		t.Run(file, func(t *testing.T) {
			data, err := ioutil.ReadFile(filepath.Join("./testdata", file))
			if err != nil {
				t.Fatal(err)
			}

			_, err = Parse(file, data)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

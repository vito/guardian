package rundmc

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/cf-guardian/specs"
)

var ErrDoesNotExist = errors.New("does not exist")

var configJson []byte

func init() {
	var err error
	if configJson, err = encodeConfigJson(); err != nil {
		panic(err) // should never happen since the above json in statically defined
	}
}

// a depot which stores containers as subdirs of a depot directory
type DirectoryDepot struct {
	Dir string
}

func (d *DirectoryDepot) Create(handle string) error {
	os.MkdirAll(filepath.Join(d.Dir, handle), 0700)

	return ioutil.WriteFile(filepath.Join(d.Dir, handle, "config.json"), configJson, 0700)
}

func (d *DirectoryDepot) Lookup(handle string) (string, error) {
	if _, err := os.Stat(filepath.Join(d.Dir, handle)); err != nil {
		return "", ErrDoesNotExist
	}

	return filepath.Join(d.Dir, handle), nil
}

func encodeConfigJson() ([]byte, error) {
	json, err := json.Marshal(specs.Spec{
		Version: "pre-draft",
		Process: specs.Process{
			Terminal: true,
			Args:     []string{"/bin/sh", "-c", `echo "Pid 1 Running"; read x`},
		},
	})

	return json, err
}

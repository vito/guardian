package rundmc

import (
	"os"
	"path/filepath"

	"github.com/pivotal-golang/lager"
)

type RootFsPreparer struct {
	OwnerUID int
	OwnerGID int
}

func (p RootFsPreparer) Make(log lager.Logger, _ bool, rootFsPath string) error {
	os.Mkdir(filepath.Join(rootFsPath, ".garden"), 0755)
	os.Chown(filepath.Join(rootFsPath, ".garden"), p.OwnerUID, p.OwnerGID)
	return nil
}

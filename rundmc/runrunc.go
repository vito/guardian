package rundmc

import (
	"encoding/json"
	"io/ioutil"
	"os/exec"

	"github.com/cf-guardian/specs"
	"github.com/cloudfoundry-incubator/garden"
	"github.com/cloudfoundry-incubator/guardian/rundmc/process_tracker"
)

//go:generate counterfeiter . ProcessTracker
type ProcessTracker interface {
	Run(id uint32, cmd *exec.Cmd, io garden.ProcessIO, tty *garden.TTYSpec, signaller process_tracker.Signaller) (garden.Process, error)
}

//go:generate counterfeiter . PidGenerator
type PidGenerator interface {
	Generate() uint32
}

// da doo
type RunRunc struct {
	Tracker      ProcessTracker
	PidGenerator PidGenerator
}

func (r RunRunc) Run(path string, io garden.ProcessIO) (garden.Process, error) {
	cmd := exec.Command("runc")
	cmd.Dir = path

	process, err := r.Tracker.Run(r.PidGenerator.Generate(), cmd, io, &garden.TTYSpec{}, nil)
	return process, err
}

func (r RunRunc) Exec(path string, spec garden.ProcessSpec, io garden.ProcessIO) (garden.Process, error) {
	tmpFile, err := ioutil.TempFile("", "guardianprocess")
	if err != nil {
		return nil, err
	}

	if err := json.NewEncoder(tmpFile).Encode(specs.Process{Args: append([]string{spec.Path}, spec.Args...)}); err != nil {
		return nil, err
	}

	cmd := exec.Command("runc", "exec", tmpFile.Name())
	cmd.Dir = path

	return r.Tracker.Run(r.PidGenerator.Generate(), cmd, io, spec.TTY, nil)
}

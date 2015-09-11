package rundmc

import (
	"io"

	"github.com/cloudfoundry-incubator/garden"
	"github.com/cloudfoundry-incubator/guardian/gardener"
)

//go:generate counterfeiter . Depot
type Depot interface {
	Create(handle string) error
	Lookup(handle string) (path string, err error)
}

//go:generate counterfeiter . ContainerRunner
type ContainerRunner interface {
	Run(path string, io garden.ProcessIO) (garden.Process, error)
	Exec(path string, spec garden.ProcessSpec, io garden.ProcessIO) (garden.Process, error)
}

type Containerizer struct {
	Depot           Depot
	ContainerRunner ContainerRunner
}

func (c *Containerizer) Create(spec gardener.DesiredContainerSpec) error {
	if err := c.Depot.Create(spec.Handle); err != nil {
		return err
	}

	path, err := c.Depot.Lookup(spec.Handle)
	if err != nil {
		return err
	}

	r, w := io.Pipe()
	c.ContainerRunner.Run(path, garden.ProcessIO{Stdout: w})
	buff := make([]byte, 1)
	r.Read(buff)

	return nil
}

func (c *Containerizer) Run(handle string, spec garden.ProcessSpec, io garden.ProcessIO) (garden.Process, error) {
	path, err := c.Depot.Lookup(handle)
	if err != nil {
		return nil, err
	}

	return c.ContainerRunner.Exec(path, spec, io)
}

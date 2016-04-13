package kawasaki

import (
	"net"
	"os"

	"github.com/pivotal-golang/lager"
)

//go:generate counterfeiter . NetnsExecer
type NetnsExecer interface {
	Exec(netnsFD *os.File, cb func() error) error
}

type configurer struct {
	hostConfigurer       HostConfigurer
	containerApplier     ContainerApplier
	instanceChainCreator InstanceChainCreator
	nsExecer             NetnsExecer
}

//go:generate counterfeiter . HostConfigurer
type HostConfigurer interface {
	Apply(logger lager.Logger, cfg NetworkConfig, netnsFD *os.File) error
	Destroy(cfg NetworkConfig) error
}

//go:generate counterfeiter . InstanceChainCreator
type InstanceChainCreator interface {
	Create(logger lager.Logger, instanceChain, bridgeName string, ip net.IP, network *net.IPNet) error
	Destroy(logger lager.Logger, instanceChain string) error
}

//go:generate counterfeiter . ContainerApplier
type ContainerApplier interface {
	Apply(logger lager.Logger, cfg NetworkConfig) error
}

func NewConfigurer(hostConfigurer HostConfigurer, containerApplier ContainerApplier, instanceChainCreator InstanceChainCreator, nsExecer NetnsExecer) *configurer {
	return &configurer{
		hostConfigurer:       hostConfigurer,
		containerApplier:     containerApplier,
		instanceChainCreator: instanceChainCreator,
		nsExecer:             nsExecer,
	}
}

func (c *configurer) Apply(log lager.Logger, cfg NetworkConfig, nsPath string) error {
	fd, err := os.Open(nsPath)
	if err != nil {
		return err
	}
	defer fd.Close()

	if err := c.hostConfigurer.Apply(log, cfg, fd); err != nil {
		return err
	}

	if err := c.instanceChainCreator.Create(log, cfg.IPTableInstance, cfg.BridgeName, cfg.ContainerIP, cfg.Subnet); err != nil {
		return err
	}

	return c.nsExecer.Exec(fd, func() error {
		return c.containerApplier.Apply(log, cfg)
	})
}

func (c *configurer) Destroy(log lager.Logger, cfg NetworkConfig) error {
	if err := c.instanceChainCreator.Destroy(log, cfg.IPTableInstance); err != nil {
		return err
	}

	return c.hostConfigurer.Destroy(cfg)
}

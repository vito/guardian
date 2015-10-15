package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"syscall"
	"time"

	"github.com/cloudfoundry-incubator/cf-debug-server"
	"github.com/cloudfoundry-incubator/cf-lager"
	"github.com/cloudfoundry-incubator/garden/server"
	"github.com/cloudfoundry-incubator/goci"
	"github.com/cloudfoundry-incubator/goci/specs"
	"github.com/cloudfoundry-incubator/guardian/gardener"
	"github.com/cloudfoundry-incubator/guardian/kawasaki"
	"github.com/cloudfoundry-incubator/guardian/kawasaki/configure"
	"github.com/cloudfoundry-incubator/guardian/kawasaki/devices"
	"github.com/cloudfoundry-incubator/guardian/kawasaki/netns"
	"github.com/cloudfoundry-incubator/guardian/kawasaki/subnets"
	"github.com/cloudfoundry-incubator/guardian/logging"
	"github.com/cloudfoundry-incubator/guardian/rundmc"
	"github.com/cloudfoundry-incubator/guardian/rundmc/depot"
	"github.com/cloudfoundry-incubator/guardian/rundmc/process_tracker"
	"github.com/cloudfoundry-incubator/guardian/rundmc/runrunc"
	"github.com/cloudfoundry/gunk/command_runner/linux_command_runner"
	"github.com/concourse/baggageclaim/volume"
	"github.com/concourse/baggageclaim/volume/driver"
	"github.com/nu7hatch/gouuid"
	"github.com/pivotal-golang/lager"
)

var PrivilegedContainerNamespaces = []specs.Namespace{
	goci.NetworkNamespace, goci.PIDNamespace, goci.UTSNamespace, goci.IPCNamespace, goci.MountNamespace,
}

var listenNetwork = flag.String(
	"listenNetwork",
	"unix",
	"how to listen on the address (unix, tcp, etc.)",
)

var listenAddr = flag.String(
	"listenAddr",
	"/tmp/garden.sock",
	"address to listen on",
)

var snapshotsPath = flag.String(
	"snapshots",
	"",
	"directory in which to store container state to persist through restarts",
)

var binPath = flag.String(
	"bin",
	"",
	"directory containing backend-specific scripts (i.e. ./create.sh)",
)

var iodaemonBin = flag.String(
	"iodaemonBin",
	"",
	"path to iodaemon binary",
)

var depotPath = flag.String(
	"depot",
	"",
	"directory in which to store containers",
)

var rootFSPath = flag.String(
	"rootfs",
	"",
	"directory of the rootfs for the containers",
)

var disableQuotas = flag.Bool(
	"disableQuotas",
	false,
	"disable disk quotas",
)

var graceTime = flag.Duration(
	"containerGraceTime",
	0,
	"time after which to destroy idle containers",
)

var portPoolStart = flag.Uint(
	"portPoolStart",
	60000,
	"start of ephemeral port range used for mapped container ports",
)

var portPoolSize = flag.Uint(
	"portPoolSize",
	5000,
	"size of port pool used for mapped container ports",
)

var networkPool = flag.String("networkPool",
	"10.254.0.0/22",
	"Pool of dynamically allocated container subnets")

var denyNetworks = flag.String(
	"denyNetworks",
	"",
	"CIDR blocks representing IPs to blacklist",
)

var allowNetworks = flag.String(
	"allowNetworks",
	"",
	"CIDR blocks representing IPs to whitelist",
)

var graphRoot = flag.String(
	"graph",
	"/var/lib/garden-docker-graph",
	"docker image graph",
)

var dockerRegistry = flag.String(
	"registry",
	"",
	///registry.IndexServerAddress(),
	"docker registry API endpoint",
)

var insecureRegistries = flag.String(
	"insecureDockerRegistryList",
	"",
	"comma-separated list of docker registries to allow connection to even if they are not secure",
)

var tag = flag.String(
	"tag",
	"",
	"server-wide identifier used for 'global' configuration, must be less than 3 character long",
)

var dropsondeOrigin = flag.String(
	"dropsondeOrigin",
	"garden-linux",
	"Origin identifier for dropsonde-emitted metrics.",
)

var dropsondeDestination = flag.String(
	"dropsondeDestination",
	"localhost:3457",
	"Destination for dropsonde-emitted metrics.",
)

var allowHostAccess = flag.Bool(
	"allowHostAccess",
	false,
	"allow network access to host",
)

var iptablesLogMethod = flag.String(
	"iptablesLogMethod",
	"kernel",
	"type of iptable logging to use, one of 'kernel' or 'nflog' (default: kernel)",
)

var mtu = flag.Int(
	"mtu",
	1500,
	"MTU size for container network interfaces")

var externalIP = flag.String(
	"externalIP",
	"",
	"IP address to use to reach container's mapped ports")

var maxContainers = flag.Uint(
	"maxContainers",
	0,
	"Maximum number of containers that can be created")

func main() {
	cf_debug_server.AddFlags(flag.CommandLine)
	cf_lager.AddFlags(flag.CommandLine)
	flag.Parse()

	logger, _ := cf_lager.New("guardian")

	if *depotPath == "" {
		missing("-depot")
	}

	if *iodaemonBin == "" {
		missing("-iodaemonBin")
	}

	resolvedRootFSPath, err := filepath.EvalSymlinks(*rootFSPath)
	if err != nil {
		panic(err)
	}

	_, networkPoolCIDR, err := net.ParseCIDR(*networkPool)
	if err != nil {
		panic(err)
	}

	backend := &gardener.Gardener{
		UidGenerator:  wireUidGenerator(),
		Starter:       wireStarter(logger),
		Networker:     wireNetworker(logger, networkPoolCIDR),
		Containerizer: wireContainerizer(logger, *depotPath, *iodaemonBin),
		VolumeCreator: wireVolumeCreator(logger, *graphRoot, resolvedRootFSPath),
		Logger:        logger,
	}

	gardenServer := server.New(*listenNetwork, *listenAddr, *graceTime, backend, logger.Session("api"))

	err = gardenServer.Start()
	if err != nil {
		logger.Fatal("failed-to-start-server", err)
	}

	signals := make(chan os.Signal, 1)

	go func() {
		<-signals
		gardenServer.Stop()
		os.Exit(0)
	}()

	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	logger.Info("started", lager.Data{
		"network": *listenNetwork,
		"addr":    *listenAddr,
	})

	select {}
}

func wireUidGenerator() gardener.UidGeneratorFunc {
	return gardener.UidGeneratorFunc(func() string { return mustStringify(uuid.NewV4()) })
}

func wireVolumeCreator(logger lager.Logger, volumeDir, defaultRootFsPath string) volume.Creator {
	volumeDriver := driver.NewBtrFSDriver(logger.Session("driver"))
	filesystem, err := volume.NewFilesystem(volumeDriver, volumeDir)

	if err != nil {
		logger.Fatal("failed-to-initialize-filesystem", err)
	}

	locker := volume.NewLockManager()
	volumeRepo := volume.NewRepository(
		logger.Session("repository"),
		filesystem,
		locker,
	)

	runner := &logging.Runner{CommandRunner: linux_command_runner.New(), Logger: logger.Session("runner")}
	provider := volume.NewGardenStrategyProvider(filesystem, runner)

	return volume.NewCreator(provider, volumeRepo, defaultRootFsPath)
}

func wireStarter(logger lager.Logger) *rundmc.Starter {
	runner := &logging.Runner{CommandRunner: linux_command_runner.New(), Logger: logger.Session("runner")}
	return rundmc.NewStarter(logger, mustOpen("/proc/cgroups"), path.Join(os.TempDir(), fmt.Sprintf("cgroups-%s", *tag)), runner)
}

func wireNetworker(log lager.Logger, networkPoolCIDR *net.IPNet) *kawasaki.Networker {
	runner := &logging.Runner{CommandRunner: linux_command_runner.New(), Logger: log.Session("runner")}

	hostCfgApplier := &configure.Host{
		Veth:   &devices.VethCreator{},
		Link:   &devices.Link{Name: "guardian"},
		Bridge: &devices.Bridge{},
		Logger: log.Session("network-host-configurer"),
	}

	containerCfgApplier := &configure.Container{
		Logger: log.Session("network-container-configurer"),
		Link:   &devices.Link{Name: "guardian"},
	}

	return kawasaki.New(
		kawasaki.NewManager(runner, "/var/run/netns"),
		kawasaki.SpecParserFunc(kawasaki.ParseSpec),
		subnets.NewPool(networkPoolCIDR),
		kawasaki.NewConfigCreator(),
		kawasaki.NewConfigApplier(
			hostCfgApplier,
			containerCfgApplier,
			&netns.Execer{},
		),
	)
}

func wireContainerizer(log lager.Logger, depotPath, iodaemonPath string) *rundmc.Containerizer {
	depot := depot.New(depotPath)

	startCheck := rundmc.StartChecker{Expect: "Pid 1 Running", Timeout: 3 * time.Second}

	runcrunner := runrunc.New(
		process_tracker.New(path.Join(os.TempDir(), fmt.Sprintf("garden-%s", *tag), "processes"), iodaemonPath, linux_command_runner.New()),
		linux_command_runner.New(),
		wireUidGenerator(),
		goci.RuncBinary("runc"),
	)

	baseBundle := goci.Bundle().
		WithNamespaces(PrivilegedContainerNamespaces...).
		WithResources(&specs.Resources{}).
		WithMounts(goci.Mount{Name: "proc", Type: "proc", Source: "proc", Destination: "/proc"}).
		WithProcess(goci.Process("/bin/sh", "-c", `echo "Pid 1 Running"; read x`))

	return rundmc.New(depot, &rundmc.BundleTemplate{Bndl: baseBundle}, runcrunner, startCheck)
}

func missing(flagName string) {
	println("missing " + flagName)
	println()
	flag.Usage()

	os.Exit(1)
}

func mustStringify(s interface{}, e error) string {
	if e != nil {
		panic(e)
	}

	return fmt.Sprintf("%s", s)
}

func mustOpen(path string) io.ReadCloser {
	if r, err := os.Open(path); err != nil {
		panic(err)
	} else {
		return r
	}
}

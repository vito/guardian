package runrunc_test

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"path"

	"github.com/cloudfoundry-incubator/garden"
	"github.com/cloudfoundry-incubator/goci"
	"github.com/cloudfoundry-incubator/guardian/rundmc/process_tracker"
	"github.com/cloudfoundry-incubator/guardian/rundmc/runrunc"
	"github.com/cloudfoundry-incubator/guardian/rundmc/runrunc/fakes"
	"github.com/cloudfoundry/gunk/command_runner/fake_command_runner"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opencontainers/runc/libcontainer/user"
	"github.com/opencontainers/specs/specs-go"
	"github.com/pivotal-golang/lager"
	"github.com/pivotal-golang/lager/lagertest"
)

var _ = Describe("Exec", func() {
	var (
		tracker       *fakes.FakeProcessTracker
		commandRunner *fake_command_runner.FakeCommandRunner
		pidGenerator  *fakes.FakeUidGenerator
		runcBinary    *fakes.FakeRuncBinary
		bundleLoader  *fakes.FakeBundleLoader
		users         *fakes.FakeUserLookupper
		mkdirer       *fakes.FakeMkdirer
		bundlePath    string
		logger        *lagertest.TestLogger
		waitWatcher   *fakes.FakeWaitWatcher

		runner       *runrunc.Execer
		execPreparer *runrunc.ExecPreparer
	)

	BeforeEach(func() {
		tracker = new(fakes.FakeProcessTracker)
		pidGenerator = new(fakes.FakeUidGenerator)
		runcBinary = new(fakes.FakeRuncBinary)
		commandRunner = fake_command_runner.New()
		logger = lagertest.NewTestLogger("test")
		bundleLoader = new(fakes.FakeBundleLoader)
		users = new(fakes.FakeUserLookupper)
		mkdirer = new(fakes.FakeMkdirer)
		waitWatcher = new(fakes.FakeWaitWatcher)

		execPreparer = runrunc.NewExecPreparer(
			bundleLoader,
			users,
			mkdirer,
		)

		var err error
		bundlePath, err = ioutil.TempDir("", "bundle")
		Expect(err).NotTo(HaveOccurred())

		runner = runrunc.NewExecer(
			execPreparer,
			runrunc.NewExecRunner(
				pidGenerator,
				runcBinary,
				tracker,
				waitWatcher,
			),
		)

		bundleLoader.LoadStub = func(path string) (*goci.Bndl, error) {
			bndl := &goci.Bndl{}
			bndl.Spec.Root.Path = rootfsPath(path)
			return bndl, nil
		}

		users.LookupReturns(&user.ExecUser{}, nil)

		runcBinary.ExecCommandStub = func(id, processJSONPath, pidFilePath string) *exec.Cmd {
			return exec.Command("funC", "exec", id, processJSONPath, "--pid-file", pidFilePath)
		}
	})

	AfterEach(func() {
		err := os.RemoveAll(bundlePath)
		Expect(err).NotTo(HaveOccurred())
	})

	It("runs exec against the injected runC binary using process tracker", func() {
		pidGenerator.GenerateReturns("another-process-guid")
		tracker.RunReturns(&process_tracker.Process{}, nil)
		ttyspec := &garden.TTYSpec{WindowSize: &garden.WindowSize{Rows: 1}}
		runner.Exec(logger, bundlePath, "some-id", garden.ProcessSpec{TTY: ttyspec}, garden.ProcessIO{Stdout: GinkgoWriter})
		Expect(tracker.RunCallCount()).To(Equal(1))

		pid, cmd, io, tty, _ := tracker.RunArgsForCall(0)
		Expect(pid).To(Equal("another-process-guid"))
		Expect(cmd.Args[:3]).To(Equal([]string{"funC", "exec", "some-id"}))
		Expect(io.Stdout).To(Equal(GinkgoWriter))
		Expect(tty).To(Equal(ttyspec))
	})

	It("creates the processes directory if it does not exist", func() {
		tracker.RunReturns(&process_tracker.Process{}, nil)
		runner.Exec(logger, bundlePath, "some-id", garden.ProcessSpec{}, garden.ProcessIO{Stdout: GinkgoWriter})
		Expect(path.Join(bundlePath, "processes")).To(BeADirectory())
	})

	Context("When creating the processes directory fails", func() {
		It("returns a helpful error", func() {
			Expect(ioutil.WriteFile(path.Join(bundlePath, "processes"), []byte(""), 0700)).To(Succeed())
			_, err := runner.Exec(logger, bundlePath, "some-id", garden.ProcessSpec{}, garden.ProcessIO{Stdout: GinkgoWriter})
			Expect(err).To(MatchError(MatchRegexp("mkdir .*: .*")))
		})
	})

	It("asks for the pid file to be placed in processes/$guid.pid", func() {
		pidGenerator.GenerateReturns("another-process-guid")
		tracker.RunReturns(&process_tracker.Process{}, nil)
		runner.Exec(logger, bundlePath, "some-id", garden.ProcessSpec{}, garden.ProcessIO{Stdout: GinkgoWriter})
		Expect(tracker.RunCallCount()).To(Equal(1))

		_, cmd, _, _, _ := tracker.RunArgsForCall(0)
		Expect(cmd.Args[4:]).To(Equal([]string{"--pid-file", path.Join(bundlePath, "/processes/another-process-guid.pid")}))
	})

	It("tells process tracker that it can find the pid-file at processes/$guid.pid", func() {
		pidGenerator.GenerateReturns("another-process-guid")
		tracker.RunReturns(&process_tracker.Process{}, nil)
		runner.Exec(logger, bundlePath, "some-id", garden.ProcessSpec{}, garden.ProcessIO{Stdout: GinkgoWriter})
		Expect(tracker.RunCallCount()).To(Equal(1))

		_, _, _, _, pidFile := tracker.RunArgsForCall(0)
		Expect(pidFile).To(Equal(path.Join(bundlePath, "/processes/another-process-guid.pid")))
	})

	It("tells runc that the process.json is in /processes/$guid.json", func() {
		pidGenerator.GenerateReturns("another-process-guid")
		tracker.RunReturns(&process_tracker.Process{}, nil)
		runner.Exec(logger, bundlePath, "some-id", garden.ProcessSpec{}, garden.ProcessIO{Stdout: GinkgoWriter})

		_, cmd, _, _, _ := tracker.RunArgsForCall(0)
		Expect(cmd.Args[3]).To(Equal(path.Join(bundlePath, "/processes/another-process-guid.json")))
	})

	Describe("process-related files", func() {
		Context("when process tracker succeeds", func() {
			var spec specs.Process
			var processJsonPath, pidFilePath string
			var fakeProcess *fakes.FakeProcess

			BeforeEach(func() {
				pidGenerator.GenerateReturns("another-process-guid")

				fakeProcess = new(fakes.FakeProcess)
				tracker.RunStub = func(_ string, cmd *exec.Cmd, _ garden.ProcessIO, _ *garden.TTYSpec, pidPath string) (garden.Process, error) {
					processJsonPath = cmd.Args[3]
					pidFilePath = pidPath

					f, err := os.Open(processJsonPath)
					Expect(err).NotTo(HaveOccurred())

					json.NewDecoder(f).Decode(&spec)
					return fakeProcess, nil
				}

				_, err := runner.Exec(logger, bundlePath, "some-id", garden.ProcessSpec{
					Path: "potato",
					Args: []string{"boom"},
				}, garden.ProcessIO{Stdout: GinkgoWriter})
				Expect(err).NotTo(HaveOccurred())
			})

			Describe("the process.json file passed to 'runc exec'", func() {
				It("is the encoded version of the config", func() {
					Expect(spec.Args).To(ConsistOf("potato", "boom"))
				})
			})

			It("defers cleanup of the process.json and pid file until the process has completed", func() {
				Eventually(waitWatcher.OnExitCallCount).Should(Equal(1))

				_, process, callback := waitWatcher.OnExitArgsForCall(0)
				Expect(process).To(Equal(fakeProcess))
				Expect(callback).To(ConsistOf(
					processJsonPath, pidFilePath,
				))
			})
		})

		Context("when process tracker fails", func() {
			It("immediately cleans up the process.json and pid file ", func() {
				var processJsonPath, pidFilePath string
				tracker.RunStub = func(_ string, cmd *exec.Cmd, _ garden.ProcessIO, _ *garden.TTYSpec, pidPath string) (garden.Process, error) {
					processJsonPath = cmd.Args[3]
					pidFilePath = pidPath

					Expect(ioutil.WriteFile(processJsonPath, []byte{}, 0700)).To(Succeed())
					Expect(ioutil.WriteFile(pidFilePath, []byte{}, 0700)).To(Succeed())

					return nil, errors.New("Boom")
				}

				_, err := runner.Exec(logger, bundlePath, "some-id", garden.ProcessSpec{}, garden.ProcessIO{Stdout: GinkgoWriter})
				Expect(err).To(MatchError(ContainSubstring("Boom")))

				Expect(processJsonPath).NotTo(BeAnExistingFile())
				Expect(pidFilePath).NotTo(BeAnExistingFile())
			})
		})
	})
})

var _ = Describe("ExecPreparer", func() {
	var (
		spec         *specs.Process
		bundleLoader *fakes.FakeBundleLoader
		users        *fakes.FakeUserLookupper
		mkdirer      *fakes.FakeMkdirer
		bundlePath   string
		logger       lager.Logger

		preparer *runrunc.ExecPreparer
	)

	BeforeEach(func() {
		logger = lagertest.NewTestLogger("test")
		bundleLoader = new(fakes.FakeBundleLoader)
		users = new(fakes.FakeUserLookupper)
		mkdirer = new(fakes.FakeMkdirer)

		var err error
		bundlePath, err = ioutil.TempDir("", "bundle")
		Expect(err).NotTo(HaveOccurred())

		bundleLoader.LoadStub = func(path string) (*goci.Bndl, error) {
			bndl := &goci.Bndl{}
			bndl.Spec.Root.Path = rootfsPath(path)
			return bndl, nil
		}

		users.LookupReturns(&user.ExecUser{}, nil)

		preparer = runrunc.NewExecPreparer(bundleLoader, users, mkdirer)
	})

	It("passes a process.json with the correct path and args", func() {
		spec, err := preparer.Prepare(logger, bundlePath, garden.ProcessSpec{Path: "to enlightenment", Args: []string{"infinity", "and beyond"}})
		Expect(err).NotTo(HaveOccurred())

		Expect(spec.Args).To(Equal([]string{"to enlightenment", "infinity", "and beyond"}))
	})

	It("sets the rlimits correctly", func() {
		ptr := func(n uint64) *uint64 { return &n }
		spec, err := preparer.Prepare(logger, bundlePath, garden.ProcessSpec{
			Limits: garden.ResourceLimits{
				As:         ptr(12),
				Core:       ptr(24),
				Cpu:        ptr(36),
				Data:       ptr(99),
				Fsize:      ptr(101),
				Locks:      ptr(111),
				Memlock:    ptr(987),
				Msgqueue:   ptr(777),
				Nice:       ptr(111),
				Nofile:     ptr(222),
				Nproc:      ptr(1234),
				Rss:        ptr(888),
				Rtprio:     ptr(254),
				Sigpending: ptr(101),
				Stack:      ptr(44),
			},
		})
		Expect(err).ToNot(HaveOccurred())

		Expect(spec.Rlimits).To(ConsistOf(
			specs.Rlimit{Type: "RLIMIT_AS", Hard: 12, Soft: 12},
			specs.Rlimit{Type: "RLIMIT_CORE", Hard: 24, Soft: 24},
			specs.Rlimit{Type: "RLIMIT_CPU", Hard: 36, Soft: 36},
			specs.Rlimit{Type: "RLIMIT_DATA", Hard: 99, Soft: 99},
			specs.Rlimit{Type: "RLIMIT_FSIZE", Hard: 101, Soft: 101},
			specs.Rlimit{Type: "RLIMIT_LOCKS", Hard: 111, Soft: 111},
			specs.Rlimit{Type: "RLIMIT_MEMLOCK", Hard: 987, Soft: 987},
			specs.Rlimit{Type: "RLIMIT_MSGQUEUE", Hard: 777, Soft: 777},
			specs.Rlimit{Type: "RLIMIT_NICE", Hard: 111, Soft: 111},
			specs.Rlimit{Type: "RLIMIT_NOFILE", Hard: 222, Soft: 222},
			specs.Rlimit{Type: "RLIMIT_NPROC", Hard: 1234, Soft: 1234},
			specs.Rlimit{Type: "RLIMIT_RSS", Hard: 888, Soft: 888},
			specs.Rlimit{Type: "RLIMIT_RTPRIO", Hard: 254, Soft: 254},
			specs.Rlimit{Type: "RLIMIT_SIGPENDING", Hard: 101, Soft: 101},
			specs.Rlimit{Type: "RLIMIT_STACK", Hard: 44, Soft: 44},
		))
	})

	It("sets Terminal to true if a TTY is configured", func() {
		spec, err := preparer.Prepare(logger, bundlePath, garden.ProcessSpec{
			TTY: &garden.TTYSpec{
				WindowSize: &garden.WindowSize{
					Columns: 80,
					Rows:    24,
				},
			},
		})
		Expect(err).ToNot(HaveOccurred())

		Expect(spec.Terminal).To(BeTrue())
	})

	Describe("passing the correct uid and gid", func() {
		Context("when the bundle can be loaded", func() {
			BeforeEach(func() {
				users.LookupReturns(&user.ExecUser{Uid: 9, Gid: 7}, nil)

				var err error
				spec, err = preparer.Prepare(logger, bundlePath, garden.ProcessSpec{User: "spiderman"})
				Expect(err).ToNot(HaveOccurred())
			})

			It("looks up the user and group IDs of the user in the right rootfs", func() {
				Expect(users.LookupCallCount()).To(Equal(1))
				actualRootfsPath, actualUserName := users.LookupArgsForCall(0)
				Expect(actualRootfsPath).To(Equal(rootfsPath(bundlePath)))
				Expect(actualUserName).To(Equal("spiderman"))
			})

			It("passes a process.json with the correct user and group ids", func() {
				Expect(spec.User).To(Equal(specs.User{UID: 9, GID: 7}))
			})
		})

		Context("when the bundle can't be loaded", func() {
			BeforeEach(func() {
				bundleLoader.LoadReturns(nil, errors.New("whoa! Hold them horses!"))
			})

			It("fails", func() {
				_, err := preparer.Prepare(logger, bundlePath,
					garden.ProcessSpec{User: "spiderman"})
				Expect(err).To(MatchError(ContainSubstring("Hold them horses")))
			})
		})

		Context("when User Lookup returns an error", func() {
			It("passes a process.json with the correct user and group ids", func() {
				users.LookupReturns(&user.ExecUser{Uid: 0, Gid: 0}, errors.New("bang"))

				_, err := preparer.Prepare(logger, bundlePath, garden.ProcessSpec{User: "spiderman"})
				Expect(err).To(MatchError(ContainSubstring("bang")))
			})
		})
	})

	Context("when the user is specified in the process spec", func() {
		Context("when the environment does not contain a USER", func() {
			It("appends a default user", func() {
				spec, err := preparer.Prepare(logger, bundlePath, garden.ProcessSpec{
					User: "spiderman",
					Env:  []string{"a=1", "b=3", "c=4", "PATH=a", "HOME=/spidermanhome"},
				})
				Expect(err).NotTo(HaveOccurred())

				Expect(spec.Env).To(ConsistOf("a=1", "b=3", "c=4", "PATH=a", "USER=spiderman", "HOME=/spidermanhome"))
			})
		})

		Context("when the environment does contain a USER", func() {
			It("appends a default user", func() {
				spec, err := preparer.Prepare(logger, bundlePath, garden.ProcessSpec{
					User: "spiderman",
					Env:  []string{"a=1", "b=3", "c=4", "PATH=a", "USER=superman"},
				})
				Expect(err).NotTo(HaveOccurred())

				Expect(spec.Env).To(Equal([]string{"a=1", "b=3", "c=4", "PATH=a", "USER=superman"}))
			})
		})
	})

	Context("when the user is not specified in the process spec", func() {
		Context("when the environment does not contain a USER", func() {
			It("passes the environment variables", func() {
				spec, err := preparer.Prepare(logger, bundlePath, garden.ProcessSpec{
					Env: []string{"a=1", "b=3", "c=4", "PATH=a"},
				})
				Expect(err).NotTo(HaveOccurred())

				Expect(spec.Env).To(Equal([]string{"a=1", "b=3", "c=4", "PATH=a", "USER=root"}))
			})
		})

		Context("when the environment already contains a USER", func() {
			It("passes the environment variables", func() {
				spec, err := preparer.Prepare(logger, bundlePath, garden.ProcessSpec{
					Env: []string{"a=1", "b=3", "c=4", "PATH=a", "USER=yo"},
				})
				Expect(err).NotTo(HaveOccurred())

				Expect(spec.Env).To(Equal([]string{"a=1", "b=3", "c=4", "PATH=a", "USER=yo"}))
			})
		})
	})

	Context("when the environment already contains a PATH", func() {
		It("passes the environment variables", func() {
			spec, err := preparer.Prepare(logger, bundlePath, garden.ProcessSpec{
				Env: []string{"a=1", "b=3", "c=4", "PATH=a"},
			})
			Expect(err).NotTo(HaveOccurred())

			Expect(spec.Env).To(Equal([]string{"a=1", "b=3", "c=4", "PATH=a", "USER=root"}))
		})
	})

	Context("when the environment does not already contain a PATH", func() {
		It("appends a default PATH for the root user", func() {
			users.LookupReturns(&user.ExecUser{Uid: 0, Gid: 0}, nil)
			spec, err := preparer.Prepare(logger, bundlePath, garden.ProcessSpec{
				Env:  []string{"a=1", "b=3", "c=4"},
				User: "root",
			})
			Expect(err).NotTo(HaveOccurred())

			Expect(spec.Env).To(Equal([]string{"a=1", "b=3", "c=4",
				"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin", "USER=root"}))
		})

		It("appends a default PATH for non-root users", func() {
			users.LookupReturns(&user.ExecUser{Uid: 1000, Gid: 1000}, nil)
			spec, err := preparer.Prepare(logger, bundlePath, garden.ProcessSpec{
				Env:  []string{"a=1", "b=3", "c=4"},
				User: "alice",
			})
			Expect(err).NotTo(HaveOccurred())

			Expect(spec.Env).To(Equal([]string{"a=1", "b=3", "c=4",
				"PATH=/usr/local/bin:/usr/bin:/bin", "USER=alice"}))
		})
	})

	Context("when the container has environment variables", func() {
		var (
			processEnv   []string
			containerEnv []string
			bndl         *goci.Bndl

			spec *specs.Process
		)

		BeforeEach(func() {
			containerEnv = []string{"ENV_CONTAINER_NAME=garden"}
			processEnv = []string{"ENV_PROCESS_ID=1"}
		})

		JustBeforeEach(func() {
			bndl = &goci.Bndl{}
			bndl.Spec.Root.Path = "/some/rootfs/path"
			bndl.Spec.Process.Env = containerEnv
			bundleLoader.LoadReturns(bndl, nil)

			var err error
			spec, err = preparer.Prepare(logger, bundlePath, garden.ProcessSpec{
				Env: processEnv,
			})
			Expect(err).NotTo(HaveOccurred())
		})

		It("appends the process vars into container vars", func() {
			envWContainer := make([]string, len(spec.Env))
			copy(envWContainer, spec.Env)

			bndl.Spec.Process.Env = []string{}
			bundleLoader.LoadReturns(bndl, nil)

			spec, err := preparer.Prepare(logger, bundlePath, garden.ProcessSpec{
				Env: processEnv,
			})
			Expect(err).NotTo(HaveOccurred())

			Expect(envWContainer).To(Equal(append(containerEnv, spec.Env...)))
		})

		Context("and the container environment contains PATH", func() {
			BeforeEach(func() {
				containerEnv = append(containerEnv, "PATH=/test")
			})

			It("should not apply the default PATH", func() {
				Expect(spec.Env).To(Equal([]string{
					"ENV_CONTAINER_NAME=garden",
					"PATH=/test",
					"ENV_PROCESS_ID=1",
					"USER=root",
				}))
			})
		})
	})

	Context("when the container has capabilities", func() {
		BeforeEach(func() {
			bndl := &goci.Bndl{}
			bndl.Spec.Process.Capabilities = []string{"foo", "bar", "baz"}
			bundleLoader.LoadReturns(bndl, nil)
		})

		It("passes them on to the process", func() {
			spec, err := preparer.Prepare(logger, bundlePath, garden.ProcessSpec{})
			Expect(err).NotTo(HaveOccurred())
			Expect(spec.Capabilities).To(Equal([]string{"foo", "bar", "baz"}))
		})
	})

	Describe("working directory", func() {
		Context("when the working directory is specified", func() {
			It("passes the correct cwd to the spec", func() {
				spec, err := preparer.Prepare(
					logger, bundlePath,
					garden.ProcessSpec{Dir: "/home/dir"},
				)
				Expect(err).NotTo(HaveOccurred())

				Expect(spec.Cwd).To(Equal("/home/dir"))
			})

			Describe("Creating the working directory", func() {
				JustBeforeEach(func() {
					users.LookupReturns(&user.ExecUser{Uid: 1012, Gid: 1013}, nil)

					_, err := preparer.Prepare(logger, bundlePath, garden.ProcessSpec{
						Dir: "/path/to/banana/dir",
					})
					Expect(err).NotTo(HaveOccurred())
				})

				Context("when the container is privileged", func() {
					It("creates the working directory", func() {
						Expect(mkdirer.MkdirAsCallCount()).To(Equal(1))
						rootfs, uid, gid, mode, recreate, dirs := mkdirer.MkdirAsArgsForCall(0)
						Expect(rootfs).To(Equal(rootfsPath(bundlePath)))
						Expect(dirs).To(ConsistOf("/path/to/banana/dir"))
						Expect(mode).To(BeNumerically("==", 0755))
						Expect(recreate).To(BeFalse())
						Expect(uid).To(BeEquivalentTo(1012))
						Expect(gid).To(BeEquivalentTo(1013))
					})
				})

				Context("when the container is unprivileged", func() {
					BeforeEach(func() {
						bundleLoader.LoadStub = func(path string) (*goci.Bndl, error) {
							bndl := &goci.Bndl{}
							bndl.Spec.Root.Path = "/rootfs/of/bundle" + path
							bndl.Spec.Linux.UIDMappings = []specs.IDMapping{{
								HostID:      1712,
								ContainerID: 1012,
								Size:        1,
							}}
							bndl.Spec.Linux.GIDMappings = []specs.IDMapping{{
								HostID:      1713,
								ContainerID: 1013,
								Size:        1,
							}}
							return bndl, nil
						}
					})

					It("creates the working directory as the mapped user", func() {
						Expect(mkdirer.MkdirAsCallCount()).To(Equal(1))
						rootfs, uid, gid, mode, recreate, dirs := mkdirer.MkdirAsArgsForCall(0)
						Expect(rootfs).To(Equal(rootfsPath(bundlePath)))
						Expect(dirs).To(ConsistOf("/path/to/banana/dir"))
						Expect(mode).To(BeEquivalentTo(0755))
						Expect(recreate).To(BeFalse())
						Expect(uid).To(BeEquivalentTo(1712))
						Expect(gid).To(BeEquivalentTo(1713))
					})
				})
			})
		})

		Context("when the working directory is not specified", func() {
			It("defaults to the user's HOME directory", func() {
				users.LookupReturns(&user.ExecUser{Home: "/the/home/dir"}, nil)

				spec, err := preparer.Prepare(
					logger, bundlePath,
					garden.ProcessSpec{Dir: ""},
				)
				Expect(err).NotTo(HaveOccurred())

				Expect(spec.Cwd).To(Equal("/the/home/dir"))
			})

			It("creates the directory", func() {
				users.LookupReturns(&user.ExecUser{Uid: 1012, Gid: 1013, Home: "/some/dir"}, nil)

				_, err := preparer.Prepare(logger, bundlePath, garden.ProcessSpec{})
				Expect(err).NotTo(HaveOccurred())

				Expect(mkdirer.MkdirAsCallCount()).To(Equal(1))
				_, _, _, _, _, dirs := mkdirer.MkdirAsArgsForCall(0)
				Expect(dirs).To(ConsistOf("/some/dir"))
			})
		})

		Context("when the working directory creation fails", func() {
			It("returns an error", func() {
				mkdirer.MkdirAsReturns(errors.New("BOOOOOM"))
				_, err := preparer.Prepare(logger, bundlePath, garden.ProcessSpec{})
				Expect(err).To(MatchError(ContainSubstring("create working directory: BOOOOOM")))
			})
		})
	})
})

var _ = Describe("WaitWatcher", func() {
	It("calls Wait only once process.Wait returns", func() {
		waiter := new(fakes.FakeWaiter)
		waitReturns := make(chan struct{})
		waiter.WaitStub = func() (int, error) {
			<-waitReturns
			return 0, nil
		}

		runner := new(fakes.FakeRunner)

		watcher := runrunc.Watcher{}
		go watcher.OnExit(lagertest.NewTestLogger("test"), waiter, runner)

		Consistently(runner.RunCallCount).ShouldNot(Equal(1))
		close(waitReturns)
		Eventually(runner.RunCallCount).Should(Equal(1))
	})
})

var _ = Describe("RemoveFiles", func() {
	It("removes all the paths", func() {
		a := tmpFile("testremovefiles")
		b := tmpFile("testremovefiles")

		runrunc.RemoveFiles([]string{a, b}).Run(lagertest.NewTestLogger("test"))
		Expect(a).NotTo(BeAnExistingFile())
		Expect(b).NotTo(BeAnExistingFile())
	})
})

func rootfsPath(bundlePath string) string {
	return "/rootfs/of/bundle" + bundlePath
}

func tmpFile(name string) string {
	tmp, err := ioutil.TempFile("", name)
	Expect(err).NotTo(HaveOccurred())
	Expect(tmp.Close()).To(Succeed())
	return tmp.Name()
}

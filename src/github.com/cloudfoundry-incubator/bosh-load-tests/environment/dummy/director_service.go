package dummy

import (
	"strconv"
	"strings"
	"time"

	bltassets "github.com/cloudfoundry-incubator/bosh-load-tests/assets"
	bltcom "github.com/cloudfoundry-incubator/bosh-load-tests/command"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
)

type DirectorService struct {
	directorMigrationCommand string
	directorStartCommand     string
	workerStartCommand       string
	assetsProvider           bltassets.Provider
	directorConfig           *DirectorConfig
	cmdRunner                boshsys.CmdRunner
	directorProcess          boshsys.Process
	workerProcesses          []boshsys.Process
	portWaiter               PortWaiter
	numWorkers               int
}

func NewDirectorService(
	directorMigrationCommand string,
	directorStartCommand string,
	workerStartCommand string,
	directorConfig *DirectorConfig,
	cmdRunner boshsys.CmdRunner,
	assetsProvider bltassets.Provider,
	portWaiter PortWaiter,
	numWorkers int,
) *DirectorService {
	return &DirectorService{
		directorMigrationCommand: directorMigrationCommand,
		directorStartCommand:     directorStartCommand,
		workerStartCommand:       workerStartCommand,
		directorConfig:           directorConfig,
		cmdRunner:                cmdRunner,
		assetsProvider:           assetsProvider,
		portWaiter:               portWaiter,
		numWorkers:               numWorkers,
	}
}

func (s *DirectorService) Start() error {
	err := s.directorConfig.Write()
	if err != nil {
		return err
	}

	migrationCommand := bltcom.CreateCommand(s.directorMigrationCommand)
	migrationCommand.Args = append(migrationCommand.Args, "-c", s.directorConfig.DirectorConfigPath())
	_, _, _, err = s.cmdRunner.RunComplexCommand(migrationCommand)
	if err != nil {
		return bosherr.WrapError(err, "running migrations")
	}

	directorCommand := bltcom.CreateCommand(s.directorStartCommand)
	directorCommand.Args = append(directorCommand.Args, "-c", s.directorConfig.DirectorConfigPath())
	s.directorProcess, err = s.cmdRunner.RunComplexCommandAsync(directorCommand)
	if err != nil {
		return bosherr.WrapError(err, "starting director")
	}

	s.directorProcess.Wait()

	err = s.portWaiter.Wait("DirectorService", "127.0.0.1", s.directorConfig.DirectorPort())
	if err != nil {
		return bosherr.WrapError(err, "Waiting for director to start up")
	}

	for i := 1; i <= s.numWorkers; i++ {
		workerStartCommand := bltcom.CreateCommand(s.workerStartCommand)
		workerStartCommand.Env["QUEUE"] = "normal"
		workerStartCommand.Args = append(workerStartCommand.Args, "-c", s.directorConfig.WorkerConfigPath(i), "-i", strconv.Itoa(i))

		workerProcess, err := s.cmdRunner.RunComplexCommandAsync(workerStartCommand)
		if err != nil {
			return bosherr.WrapError(err, "starting worker")
		}
		workerProcess.Wait()
		s.workerProcesses = append(s.workerProcesses, workerProcess)

		if err != nil {
			return bosherr.WrapError(err, "Waiting for worker to start up")
		}
	}

	return s.waitForWorkersToStart()
}

func (s *DirectorService) Stop() {
	for _, process := range s.workerProcesses {
		process.TerminateNicely(5 * time.Second)
	}
	s.directorProcess.TerminateNicely(5 * time.Second)

	for _, worker := range s.workerProcesses {
		worker.TerminateNicely(5 * time.Second)
	}
}

func (s *DirectorService) waitForWorkersToStart() error {
	cmd := boshsys.Command{
		Name: "bash",
		Args: []string{"-c", "ps ax | grep bosh-director/bin/bosh-director-worker | grep -v grep | wc -l"},
	}

	for i := 0; i < 30; i++ {
		stdout, _, _, _ := s.cmdRunner.RunComplexCommand(cmd)
		if strings.TrimSpace(stdout) == strconv.Itoa(s.numWorkers) {
			return nil
		}
		time.Sleep(1 * time.Second)
	}

	return bosherr.Error("Timed out waiting for workers to start")
}

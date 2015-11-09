package dummy

import (
	"time"

	bltassets "github.com/cloudfoundry-incubator/bosh-load-tests/assets"
	bltcom "github.com/cloudfoundry-incubator/bosh-load-tests/command"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
)

type NginxService struct {
	nginxStartCommand string
	directorPort      int
	nginxPort         int
	cmdRunner         boshsys.CmdRunner
	process           boshsys.Process
	assetsProvider    bltassets.Provider
	portWaiter        PortWaiter
}

func NewNginxService(
	nginxStartCommand string,
	directorPort int,
	nginxPort int,
	cmdRunner boshsys.CmdRunner,
	assetsProvider bltassets.Provider,
	portWaiter PortWaiter,
) *NginxService {
	return &NginxService{
		nginxStartCommand: nginxStartCommand,
		directorPort:      directorPort,
		nginxPort:         nginxPort,
		cmdRunner:         cmdRunner,
		assetsProvider:    assetsProvider,
		portWaiter:        portWaiter,
	}
}

func (s *NginxService) Start() error {
	nginxStartCommand := bltcom.CreateCommand(s.nginxStartCommand)
	configPath, err := s.assetsProvider.FullPath("nginx.yml")
	if err != nil {
		return bosherr.WrapError(err, "Getting nginx config path")
	}

	nginxStartCommand.Args = append(nginxStartCommand.Args, "-c", configPath)

	s.process, err = s.cmdRunner.RunComplexCommandAsync(nginxStartCommand)
	if err != nil {
		return bosherr.WrapError(err, "starting nginx")
	}

	s.process.Wait()
	return s.portWaiter.Wait("NginxService", "127.0.0.1", s.nginxPort)
}

func (s *NginxService) Stop() {
	s.process.TerminateNicely(5 * time.Second)
}

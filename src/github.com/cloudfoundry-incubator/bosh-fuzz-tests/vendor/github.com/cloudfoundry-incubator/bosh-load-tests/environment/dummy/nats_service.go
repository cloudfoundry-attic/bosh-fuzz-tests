package dummy

import (
	"bytes"
	"text/template"
	"time"

	bltassets "github.com/cloudfoundry-incubator/bosh-load-tests/assets"
	bltcom "github.com/cloudfoundry-incubator/bosh-load-tests/command"
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
)

type NatsServerOptions struct {
	AssetsPath string
	Port       int
}

type NatsService struct {
	options          NatsServerOptions
	natsStartCommand string
	cmdRunner        boshsys.CmdRunner
	process          boshsys.Process
	assetsProvider   bltassets.Provider
	fs               boshsys.FileSystem
	portWaiter       PortWaiter
}

func NewNatsService(
	options NatsServerOptions,
	natsStartCommand string,
	cmdRunner boshsys.CmdRunner,
	assetsProvider bltassets.Provider,
	fs boshsys.FileSystem,
	portWaiter PortWaiter,
) *NatsService {
	return &NatsService{
		options:          options,
		natsStartCommand: natsStartCommand,
		cmdRunner:        cmdRunner,
		assetsProvider:   assetsProvider,
		fs:               fs,
		portWaiter:       portWaiter,
	}
}

func (s *NatsService) Start() error {
	natsTemplatePath, err := s.assetsProvider.FullPath("nats/nats.template")
	if err != nil {
		return err
	}

	configPath, err := s.renderConfig(natsTemplatePath)
	if err != nil {
		return err
	}

	natsStartCommand := bltcom.CreateCommand(s.natsStartCommand)
	natsStartCommand.Args = append(natsStartCommand.Args, "-c", configPath, "-T", "-D", "-V")

	s.process, err = s.cmdRunner.RunComplexCommandAsync(natsStartCommand)
	if err != nil {
		return bosherr.WrapError(err, "starting nats server")
	}

	s.process.Wait()

	return s.portWaiter.Wait("NatsService", "127.0.0.1", s.options.Port)
}

func (s *NatsService) Stop() {
	s.process.TerminateNicely(5 * time.Second)
}

func (s *NatsService) renderConfig(natsTemplatePath string) (string, error) {
	configTemplate := template.Must(template.ParseFiles(natsTemplatePath))
	buffer := bytes.NewBuffer([]byte{})

	err := configTemplate.Execute(buffer, s.options)
	if err != nil {
		return "", err
	}

	renderedConfig, err := s.fs.TempFile("nats-config")
	if err != nil {
		return "", err
	}

	err = s.fs.WriteFile(renderedConfig.Name(), buffer.Bytes())
	if err != nil {
		return "", err
	}

	return renderedConfig.Name(), nil
}

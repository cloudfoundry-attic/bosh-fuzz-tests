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

type ConfigServerOptions struct {
	AssetsPath string
	Port       int
	Store      string
}

type ConfigServerService struct {
	options        ConfigServerOptions
	startCommand   string
	cmdRunner      boshsys.CmdRunner
	process        boshsys.Process
	assetsProvider bltassets.Provider
	fs             boshsys.FileSystem
	portWaiter     PortWaiter
}

func NewConfigServerService(
	startCommand string,
	options ConfigServerOptions,
	cmdRunner boshsys.CmdRunner,
	assetsProvider bltassets.Provider,
	fs boshsys.FileSystem,
	portWaiter PortWaiter,
) *ConfigServerService {
	return &ConfigServerService{
		startCommand:   startCommand,
		options:        options,
		cmdRunner:      cmdRunner,
		assetsProvider: assetsProvider,
		fs:             fs,
		portWaiter:     portWaiter,
	}
}

func (u *ConfigServerService) Start() error {
	configTemplatePath, err := u.assetsProvider.FullPath("config_server/config.template")
	if err != nil {
		return err
	}

	configPath, err := u.renderConfig(configTemplatePath)
	if err != nil {
		return err
	}

	command := bltcom.CreateCommand(u.startCommand)
	command.Args = append(command.Args, configPath)

	u.process, err = u.cmdRunner.RunComplexCommandAsync(command)
	if err != nil {
		return bosherr.WrapError(err, "starting config server")
	}

	return u.portWaiter.Wait("ConfigServerService", "localhost", u.options.Port)
}

func (u *ConfigServerService) Stop() {
	u.process.TerminateNicely(5 * time.Second)
}

func (u *ConfigServerService) renderConfig(configTemplatePath string) (string, error) {
	configTemplate := template.Must(template.ParseFiles(configTemplatePath))
	buffer := bytes.NewBuffer([]byte{})

	err := configTemplate.Execute(buffer, u.options)
	if err != nil {
		return "", err
	}

	renderedConfig, err := u.fs.TempFile("config")
	if err != nil {
		return "", err
	}

	err = u.fs.WriteFile(renderedConfig.Name(), buffer.Bytes())
	if err != nil {
		return "", err
	}

	return renderedConfig.Name(), nil
}

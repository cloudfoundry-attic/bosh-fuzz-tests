package dummy

import (
	"errors"
	"os"
	"path/filepath"
	"time"

	bltassets "github.com/cloudfoundry-incubator/bosh-load-tests/assets"
	bltconfig "github.com/cloudfoundry-incubator/bosh-load-tests/config"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
)

type dummy struct {
	workingDir          string
	database            Database
	directorService     *DirectorService
	nginxService        *NginxService
	uaaService          *UAAService
	configServerService *ConfigServerService
	natsService         *NatsService
	config              *bltconfig.Config
	fs                  boshsys.FileSystem
	cmdRunner           boshsys.CmdRunner
	assetsProvider      bltassets.Provider
	logger              boshlog.Logger
}

func NewDummy(
	config *bltconfig.Config,
	fs boshsys.FileSystem,
	cmdRunner boshsys.CmdRunner,
	assetsProvider bltassets.Provider,
	logger boshlog.Logger,

) *dummy {
	return &dummy{
		config:         config,
		fs:             fs,
		cmdRunner:      cmdRunner,
		assetsProvider: assetsProvider,
		logger:         logger,
	}
}

func (d *dummy) Setup() error {
	var err error
	d.workingDir, err = d.fs.TempDir("dummy-working-dir")
	if err != nil {
		return err
	}

	if "mysql" == os.Getenv("DB") {
		d.database = NewMysqlDatabase("test", d.cmdRunner)
	} else if "postgresql" == os.Getenv("DB") {
		d.database = NewPostgresqlDatabase("test", d.cmdRunner)
	} else {
		return errors.New("Unrecognized database server. Please use the DB environment variable to set database server to postgresql or mysql.")
	}

	err = d.database.Create()
	if err != nil {
		return err
	}

	portWaiter := NewPortWaiter(30, 1*time.Second)

	d.natsService = NewNatsService(d.config.NatsStartCommand, 65010, d.cmdRunner, portWaiter)
	err = d.natsService.Start()
	if err != nil {
		return err
	}

	d.nginxService = NewNginxService(d.config.NginxStartCommand, 65001, 65002, d.cmdRunner, d.assetsProvider, portWaiter)
	err = d.nginxService.Start()
	if err != nil {
		return err
	}

	if d.config.UAAConfig.Enabled {
		uaaOptions := UAAServiceOptions{
			AssetsPath:            d.config.AssetsPath,
			TomcatPath:            d.config.UAAConfig.TomcatPath,
			UaaHttpPort:           65003,
			UaaServerPort:         65004,
			UaaAccessLogDirectory: filepath.Join(d.workingDir, "UaaAccessLogDirectory"),
		}

		d.uaaService = NewUAAService(uaaOptions, d.cmdRunner, d.assetsProvider, d.fs, d.logger)
		err = d.uaaService.Start()
		if err != nil {
			return err
		}
	}

	if d.config.ConfigServerConfig.Enabled {
		if !d.config.UAAConfig.Enabled {
			return errors.New("Config server requires UAA")
		}

		configServerOptions := ConfigServerOptions{
			AssetsPath: d.config.AssetsPath,
			Port:       65005,
			Store:      "memory",
		}

		d.configServerService = NewConfigServerService(d.config.ConfigServerConfig.ConfigServerStartCommand, configServerOptions, d.cmdRunner, d.assetsProvider, d.fs, portWaiter)
		err = d.configServerService.Start()
		if err != nil {
			return err
		}
	}

	directorOptions := DirectorOptions{
		Port:                  65001,
		DatabaseName:          d.database.Name(),
		DatabaseServer:        d.database.Server(),
		DatabaseUser:          d.database.User(),
		DatabasePassword:      d.database.Password(),
		DatabasePort:          d.database.Port(),
		BaseDir:               d.workingDir,
		DummyCPIPath:          d.config.DummyCPIPath,
		RubyVersion:           d.config.RubyVersion,
		VerifyMultidigestPath: d.config.VerifyMultidigest,
		UAAEnabled:            d.config.UAAConfig.Enabled,
		ConfigServerEnabled:   d.config.ConfigServerConfig.Enabled,
		AssetsPath:            d.config.AssetsPath,
	}

	directorConfig := NewDirectorConfig(directorOptions, d.fs, d.assetsProvider, d.config.NumberOfWorkers)
	d.directorService = NewDirectorService(
		d.config.DirectorMigrationCommand,
		d.config.DirectorStartCommand,
		d.config.WorkerStartCommand,
		directorConfig,
		d.cmdRunner,
		d.assetsProvider,
		portWaiter,
		d.config.NumberOfWorkers,
	)

	err = d.directorService.Start()
	if err != nil {
		return err
	}

	return nil
}

func (d *dummy) Shutdown() error {
	return nil
}

func (d *dummy) DirectorURL() string {
	return "http://localhost:65002"
}

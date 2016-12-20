package dummy

import (
	"time"

	bltassets "github.com/cloudfoundry-incubator/bosh-load-tests/assets"
	bltconfig "github.com/cloudfoundry-incubator/bosh-load-tests/config"

	"errors"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
	"os"
)

type dummy struct {
	workingDir      string
	database        Database
	directorService *DirectorService
	nginxService    *NginxService
	natsService     *NatsService
	config          *bltconfig.Config
	fs              boshsys.FileSystem
	cmdRunner       boshsys.CmdRunner
	assetsProvider  bltassets.Provider
}

func NewDummy(
	config *bltconfig.Config,
	fs boshsys.FileSystem,
	cmdRunner boshsys.CmdRunner,
	assetsProvider bltassets.Provider,
) *dummy {
	return &dummy{
		config:         config,
		fs:             fs,
		cmdRunner:      cmdRunner,
		assetsProvider: assetsProvider,
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

	directorOptions := DirectorOptions{
		Port:             65001,
		DatabaseName:     d.database.Name(),
		DatabaseServer:   d.database.Server(),
		DatabaseUser:     d.database.User(),
		DatabasePassword: d.database.Password(),
		DatabasePort:     d.database.Port(),
	}

	directorConfig := NewDirectorConfig(directorOptions, d.workingDir, d.fs, d.assetsProvider, d.config.NumberOfWorkers, d.config.DummyCPIPath)
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

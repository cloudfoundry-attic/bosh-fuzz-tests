package dummy

import (
	"errors"
	"strings"

	boshsys "github.com/cloudfoundry/bosh-utils/system"
	boshuuid "github.com/cloudfoundry/bosh-utils/uuid"
)

type Database struct {
	name      string
	cmdRunner boshsys.CmdRunner
}

func NewDatabase(name string, cmdRunner boshsys.CmdRunner) *Database {
	return &Database{
		cmdRunner: cmdRunner,
	}
}

func (d *Database) Name() string {
	return d.name
}

func (d *Database) Create() error {
	uuid, err := boshuuid.NewGenerator().Generate()
	if err != nil {
		return err
	}
	d.name = strings.Join([]string{"bosh", uuid}, "-")

	d.Drop()
	_, _, _, err = d.cmdRunner.RunCommand("psql", "-U", "postgres", "-c", "create database \""+d.name+"\";")
	if err != nil {
		return err
	}
	return nil
}

func (d *Database) Drop() error {
	if d.name == "" {
		return errors.New("Need to create database first")
	}

	_, _, _, err := d.cmdRunner.RunCommand("psql", "-U", "postgres", "-c", "drop database \""+d.name+"\";")
	if err != nil {
		return err
	}
	return nil
}

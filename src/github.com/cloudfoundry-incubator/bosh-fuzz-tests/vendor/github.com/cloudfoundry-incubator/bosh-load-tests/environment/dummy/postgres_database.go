package dummy

import (
	"errors"
	"strings"

	boshsys "github.com/cloudfoundry/bosh-utils/system"
	boshuuid "github.com/cloudfoundry/bosh-utils/uuid"
)

type PostgresqlDatabase struct {
	name      string
	cmdRunner boshsys.CmdRunner
}

func NewPostgresqlDatabase(name string, cmdRunner boshsys.CmdRunner) *PostgresqlDatabase {
	return &PostgresqlDatabase{
		name: name,
		cmdRunner: cmdRunner,
	}
}

func (d *PostgresqlDatabase) Server() string {
	return "postgresql"
}

func (d *PostgresqlDatabase) User() string {
	return "postgres"
}

func (d *PostgresqlDatabase) Password() string {
	return "password"
}

func (d *PostgresqlDatabase) Port() int {
	return 5432
}

func (p *PostgresqlDatabase) Name() string {
	return p.name
}

func (p *PostgresqlDatabase) Create() error {
	uuid, err := boshuuid.NewGenerator().Generate()
	if err != nil {
		return err
	}
	p.name = strings.Join([]string{"bosh", uuid}, "-")

	p.Drop()
	_, _, _, err = p.cmdRunner.RunCommand("psql", "-U", p.User(), "-c", "create database \""+p.name+"\";")
	if err != nil {
		return err
	}
	return nil
}

func (p *PostgresqlDatabase) Drop() error {
	if p.name == "" {
		return errors.New("Need to create database first")
	}

	_, _, _, err := p.cmdRunner.RunCommand("psql", "-U", p.User(), "-c", "drop database \""+p.name+"\";")
	if err != nil {
		return err
	}
	return nil
}

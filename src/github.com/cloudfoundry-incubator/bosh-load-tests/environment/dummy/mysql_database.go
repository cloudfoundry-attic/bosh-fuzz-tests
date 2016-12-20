package dummy

import (
	"errors"
	"strings"

	boshsys "github.com/cloudfoundry/bosh-utils/system"
	boshuuid "github.com/cloudfoundry/bosh-utils/uuid"
)

type MysqlDatabase struct {
	name      string
	cmdRunner boshsys.CmdRunner
}

func NewMysqlDatabase(name string, cmdRunner boshsys.CmdRunner) *MysqlDatabase {
	return &MysqlDatabase{
		cmdRunner: cmdRunner,
	}
}

func (m *MysqlDatabase) Server() string {
	return "mysql2"
}

func (m *MysqlDatabase) User() string {
	return "root"
}

func (m *MysqlDatabase) Password() string {
	return "password"
}

func (m *MysqlDatabase) Port() int {
	return 3306
}

func (m *MysqlDatabase) Name() string {
	return m.name
}

func (m *MysqlDatabase) Create() error {
	uuid, err := boshuuid.NewGenerator().Generate()
	if err != nil {
		return err
	}
	m.name = strings.Join([]string{"bosh", uuid}, "-")

	m.Drop()
	_, _, _, err = m.cmdRunner.RunCommand("mysql", "--user="+m.User(), "--password="+m.Password(), "-e", "create database `"+m.name+"`")
	if err != nil {
		return err
	}
	return nil
}

func (m *MysqlDatabase) Drop() error {
	if m.name == "" {
		return errors.New("Need to create database first")
	}
	_, _, _, err := m.cmdRunner.RunCommand("mysql", "--user="+m.User(), "--password="+m.Password(), "-e", "drop database `"+m.name+"`")
	if err != nil {
		return err
	}
	return nil
}

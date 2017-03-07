package command

import (
	"strings"

	boshsys "github.com/cloudfoundry/bosh-utils/system"
)

func CreateCommand(command string) boshsys.Command {
	cmdParts := strings.Split(command, " ")
	args := []string{}
	env := map[string]string{}
	var name string

	for i := 0; i < len(cmdParts); i++ {
		if strings.Contains(cmdParts[i], "=") {
			envPair := strings.Split(cmdParts[i], "=")
			env[envPair[0]] = envPair[1]
			continue
		}

		if name == "" {
			name = cmdParts[i]
			continue
		}

		args = append(args, cmdParts[i])
	}

	return boshsys.Command{
		Name: name,
		Args: args,
		Env:  env,
	}
}

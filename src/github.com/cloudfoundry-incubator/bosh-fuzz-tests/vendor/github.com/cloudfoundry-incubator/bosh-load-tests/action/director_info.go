package action

import (
	"encoding/json"

	bltclirunner "github.com/cloudfoundry-incubator/bosh-load-tests/action/clirunner"
)

type DirectorInfo struct {
	UUID string
	URL  string
	Name string
}

func NewDirectorInfo(directorURL string, cliRunnerFactory bltclirunner.Factory) (DirectorInfo, error) {
	cliRunner, err := cliRunnerFactory.Create("bosh")
	if nil != err {
		return DirectorInfo{}, err
	}

	cliRunner.SetEnv(directorURL)

	output, err := cliRunner.RunWithOutput("env", "--json")
	if err != nil {
		return DirectorInfo{}, err
	}

	var outputStruct Output
	json.Unmarshal([]byte(output), &outputStruct)

	var name, uuid string
	for _, row := range outputStruct.Tables[0].Rows {
		switch row["col_0"] {
		case "Name":
			name = row["col_1"].(string)
		case "UUID":
			uuid = row["col_1"].(string)
		default:
			continue
		}
	}

	return DirectorInfo{
		UUID: uuid,
		URL:  directorURL,
		Name: name,
	}, nil
}

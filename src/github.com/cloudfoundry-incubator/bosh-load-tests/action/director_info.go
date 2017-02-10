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

type Table struct {
	Rows [][]string `json:"Rows"`
}

type Output struct {
	Tables []Table `json:"Tables"`
}

func NewDirectorInfo(directorURL string, cliRunnerFactory bltclirunner.Factory) (DirectorInfo, error) {
	cliRunner := cliRunnerFactory.Create()

	cliRunner.SetEnv(directorURL)

	output, err := cliRunner.RunWithOutput("env", "--json")
	if err != nil {
		return DirectorInfo{}, err
	}

	var outputStruct Output
	json.Unmarshal([]byte(output), &outputStruct)

	var name, uuid string
	for _, row := range outputStruct.Tables[0].Rows {
		switch row[0] {
		case "Name":
			name = row[1]
		case "UUID":
			uuid = row[1]
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

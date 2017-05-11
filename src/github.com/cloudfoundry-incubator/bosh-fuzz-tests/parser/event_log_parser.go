package parser

import (
	"encoding/json"
	"strings"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
)

type EventOutput struct {
	Tables []Table
}

type Table struct {
	Rows Events
}

type Events []Event

type Event struct {
	Action     string `json:"action"`
	Context    string `json:"context"`
	Deployment string `json:"deployment"`
	Error      string `json:"error"`
	Id         string `json:"id"`
	Instance   string `json:"instance"`
	ObjectName string `json:"object_name"`
	ObjectType string `json:"object_type"`
	TaskId     string `json:"task_id"`
	Time       string `json:"time"`
	User       string `json:"user"`
}

func ParseEventLog(eventLog string) (Events, error) {
	eventOutput := EventOutput{}
	err := json.Unmarshal([]byte(eventLog), &eventOutput)
	if err != nil {
		return nil, err
	}
	return eventOutput.Tables[0].Rows, nil
}

func (e Events) FindById(id string) (Event, error) {
	var expectedId string
	if strings.HasPrefix(id, "/") {
		expectedId = id
	} else {
		expectedId = "/TestDirector/foo-deployment/" + id
	}

	for _, event := range e {
		if event.ObjectName == expectedId {
			return event, nil
		}
	}
	return Event{}, bosherr.Errorf("Event with ObjectId '%s' not found", expectedId)
}

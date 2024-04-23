package events

import (
	"encoding/json"
	"github.com/google/uuid"
)

type PlanTriggered struct {
	ProjectId    uuid.UUID
	PlanId       uuid.UUID
	RepoPath     string
	CommitSha    string
	Variables    map[string]string
	EnvVariables map[string]string
}

func (i PlanTriggered) MarshalBinary() ([]byte, error) {
	return json.Marshal(i)
}

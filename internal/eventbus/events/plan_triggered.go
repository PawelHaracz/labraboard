package events

import (
	"encoding/json"
	"github.com/google/uuid"
	"labraboard/internal/models"
)

type Commit struct {
	Type models.CommitType
	Name string
}

type PlanTriggered struct {
	ProjectId    uuid.UUID
	PlanId       uuid.UUID
	RepoPath     string
	Commit       Commit
	Variables    map[string]string
	EnvVariables map[string]string
}

func (i PlanTriggered) MarshalBinary() ([]byte, error) {
	return json.Marshal(i)
}

package aggregates

import (
	"github.com/google/uuid"
	tfjson "github.com/hashicorp/terraform-json"
	"time"
)

type TerraformState struct {
	projectId uuid.UUID
	state     []byte
	CreatedOn time.Time
	ModifyOn  time.Time
}

func NewTerraformState(projectId uuid.UUID, state []byte) (*TerraformState, error) {
	utc := time.Now().UTC()
	return &TerraformState{projectId: projectId,
		state:     state,
		CreatedOn: utc,
		ModifyOn:  utc,
	}, nil
}

func (s *TerraformState) GetState() (*tfjson.State, error) {
	state := tfjson.State{}
	err := state.UnmarshalJSON(s.state)
	if err != nil {
		return nil, err
	}
	return &state, nil
}

func (s *TerraformState) GetByteState() []byte {
	return s.state
}

func (s *TerraformState) SetState(state *[]byte) {
	utc := time.Now().UTC()
	s.state = make([]byte, len(*state))
	copy(s.state, *state)
	s.ModifyOn = utc
}

func (s *TerraformState) GetID() uuid.UUID {
	return s.projectId
}

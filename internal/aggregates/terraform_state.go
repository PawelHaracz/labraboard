package aggregates

import (
	"bytes"
	"encoding/json"
	"github.com/google/uuid"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/pkg/errors"
	"time"
)

type TerraformState struct {
	projectId uuid.UUID
	state     []byte
	CreatedOn time.Time
	ModifyOn  time.Time
}

type StatFile struct {
	CheckResults     []interface{}          `json:"check_results"`
	Outputs          map[string]interface{} `json:"outputs"`
	Resources        []interface{}          `json:"resources"`
	TerraformVersion string                 `json:"terraform_version"`
	Version          int                    `json:"version"`
	Lineage          uuid.UUID              `json:"lineage"`
	Serial           int64                  `json:"serial"`
	FormatVersion    string                 `json:"format_version"`
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

func (s *TerraformState) Serialize(state *StatFile) ([]byte, error) {
	return json.Marshal(state)
}

func (s *TerraformState) Deserialize() (*StatFile, error) {
	var state StatFile
	r := bytes.NewReader(s.state)
	if err := json.NewDecoder(r).Decode(&state); err != nil {
		return nil, errors.Wrap(err, "cannot decode state")
	}
	state.FormatVersion = "1.0"
	return &state, nil
}

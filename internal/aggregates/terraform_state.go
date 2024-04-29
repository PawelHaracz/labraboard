package aggregates

import (
	"bytes"
	"encoding/json"
	"fmt"
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
	lock      []byte
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

type LockInfo struct {
	// "Created": "2024-02-05T20:04:43.120857Z",
	Created string `json:"Created"`
	// "ID": "5b64957f-e4d3-8820-77a2-913e4a8a10bd",
	ID string `json:"ID"`
	// "Info": "",
	Info string `json:"Info"`
	// "Operation": "OperationTypePlan",
	Operation string `json:"Operation"`
	// "Path": "",
	Path string `json:"Path"`
	// "Version": "1.7.2",
	Version string `json:"Version"`
	// "Who": "nhruby@newhope.local"
	Who string `json:"Who"`
}

func NewTerraformState(projectId uuid.UUID, state []byte, on time.Time, modifyOn time.Time, lock []byte) (*TerraformState, error) {
	return &TerraformState{projectId: projectId,
		state:     state,
		CreatedOn: on,
		ModifyOn:  modifyOn,
		lock:      lock,
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
func (s *TerraformState) GetByteLock() []byte {
	return s.lock
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
	if len(s.state) == 0 {
		return nil, nil
	}
	r := bytes.NewReader(s.state)
	if err := json.NewDecoder(r).Decode(&state); err != nil {
		return nil, errors.Wrap(err, "cannot decode state")
	}
	state.FormatVersion = "1.0"
	return &state, nil
}

func (s *TerraformState) GetLockInfo() (*LockInfo, error) {
	if len(s.lock) == 0 {
		return nil, nil
	}
	var storedLock LockInfo
	err := json.Unmarshal(s.lock, &storedLock)
	if err != nil {
		return nil, errors.Wrap(err, "cannot decode lock")
	}
	return &storedLock, nil
}

func (s *TerraformState) SetLockInfo(lock *LockInfo) error {
	if lock == nil {
		s.lock = make([]byte, 0)
		return nil
	}
	b, err := json.Marshal(lock)
	if err != nil {
		return errors.Wrap(err, "cannot encode lock")
	}
	s.lock = make([]byte, len(b))
	copy(s.lock, b)
	return nil
}

func (s *TerraformState) LeaseLock(reqLock *LockInfo) error {
	storedLock, err := s.GetLockInfo()
	if err != nil {
		return errors.Wrap(err, "cannot get lock")
	}
	if storedLock != nil && reqLock.ID != storedLock.ID {
		return errors.New(fmt.Sprintf("lock %s is already taken", reqLock.ID))
	}
	s.lock = nil
	return nil
}

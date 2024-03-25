package memory

import (
	"errors"
	"github.com/google/uuid"
	"labraboard/internal/aggregates"
	vo "labraboard/internal/valueobjects"
	"testing"
)

func TestMemory_GetIac(t *testing.T) {
	type testCase struct {
		name        string
		id          uuid.UUID
		expectedErr error
	}

	// Create a fake customer to add to repositories
	preparedId, _ := uuid.Parse("f47ac10b-58cc-0372-8567-0e02b2c3d479")
	iac, err := aggregates.NewIac(preparedId, vo.Terraform)
	if err != nil {
		t.Fatal(err)
	}
	id := iac.GetID()
	// Create the repo to use, and add some test Data to it for testing
	// Skip Factory for this
	repo := Repository{
		iacs: map[uuid.UUID]*aggregates.Iac{
			id: iac,
		},
	}

	testCases := []testCase{
		{
			id:          id,
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			_, err := repo.Get(tc.id)
			if !errors.Is(err, tc.expectedErr) {
				t.Errorf("Expected error %v, got %v", tc.expectedErr, err)
			}
		})
	}
}

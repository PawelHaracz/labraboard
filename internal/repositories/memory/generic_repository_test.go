package memory

import (
	"errors"
	"github.com/google/uuid"
	"golang.org/x/net/context"
	"labraboard/internal/aggregates"
	vo "labraboard/internal/valueobjects"
	"testing"
)

func TestGenericMemory_GetIac(t *testing.T) {
	type testCase struct {
		name        string
		id          uuid.UUID
		expectedErr error
	}

	// Create a fake customer to add to repositories
	preparedId, _ := uuid.Parse("f47ac10b-58cc-0372-8567-0e02b2c3d479")
	iac, err := aggregates.NewIac(preparedId, vo.Terraform, make([]*vo.Plans, 0), make([]*vo.IaCEnv, 0), nil, make([]*vo.IaCVariable, 0))
	if err != nil {
		t.Fatal(err)
	}
	id := iac.GetID()

	repo := NewGenericRepository[*aggregates.Iac]()

	testCases := []testCase{
		{
			id:          id,
			expectedErr: nil,
		},
	}
	repo.Add(iac, context.TODO())
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			_, err := repo.Get(tc.id, context.TODO())
			if !errors.Is(err, tc.expectedErr) {
				t.Errorf("Expected error %v, got %v", tc.expectedErr, err)
			}
		})
	}
}

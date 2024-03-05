package iac

import (
	"github.com/google/uuid"
	"testing"
)

func TestNewTofuIacService(t *testing.T) {
	iac, err := NewTofuIacService("/Users/pawelharacz/src/PoC/tf-example")
	if err != nil {
		t.Errorf("error: %v", err)
	}

	p, err := iac.Plan(uuid.New())
	if err != nil {
		t.Errorf("error: %v", err)
	}

	if p.Type != Tofu {
		t.Errorf("tofu type is not set")
	}

}

package iac

import (
	"github.com/google/uuid"
	"testing"
)

func TestNewTofuIacService(t *testing.T) {
	iac, err := NewTofuIacService("/Users/pawelharacz/src/PoC/tf-example", true)
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	p, err := iac.Plan(uuid.New())
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	if p.Type != Tofu {
		t.Fatalf("tofu type is not set")
	}

	if p.plan == nil {
		t.Fatalf("tofu type is not set")
	}

	add, update, deleted := p.plan.GetChanges()
	if add != 5 {
		t.Fatalf("Add not equla expetect value")
	}
	if update != 0 {
		t.Fatalf("Update not equla expetect value")
	}
	if deleted != 0 {
		t.Fatalf("Delete not equla expetect value")
	}
}

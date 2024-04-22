package iac

import (
	"github.com/go-playground/assert/v2"
	"testing"
)

// TestNewTofuIacService tests the NewTofuIacService function you have to have permissions for perform terraform on cloud
func TestNewTofuIacService(t *testing.T) {
	iac, err := NewTofuIacService("/Users/pawelharacz/src/PoC/tf-example")
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	p, err := iac.Plan(nil, nil)
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	assert.NotEqual(t, p.GetPlan(), nil)
}

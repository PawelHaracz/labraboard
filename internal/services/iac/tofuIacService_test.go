package iac

import (
	"github.com/go-playground/assert/v2"
	"golang.org/x/net/context"
	"testing"
)

// TestNewTofuIacService tests the NewTofuIacService function you have to have permissions for perform terraform on cloud
func TestNewTofuIacService(t *testing.T) {
	t.SkipNow()
	ctx := context.TODO()
	iac, err := NewTofuIacService("/Users/pawelharacz/src/PoC/tf-example", ctx)
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	p, err := iac.Plan(nil, nil, ctx)
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	planJson, planRaw := p.GetPlan()
	assert.NotEqual(t, planJson, nil)
	assert.NotEqual(t, planRaw, nil)
}

package client

import (
	"context"
	"testing"
)

func TestGetPolicies(t *testing.T) {
	client := NewClient("http://localhost:3001", WithAPIKey("GSHMG1A56JWNRX29YXE1IJQ0064QCXRL"))
	ctx := context.Background()
	policies, err := client.GetPolicies(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(policies) != 0 {
		t.Errorf("got %d policies, wanted %d policies", len(policies), 0)
	}
}

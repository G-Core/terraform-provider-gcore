package gcore

import (
	"testing"

	"github.com/G-Core/gcorelabscloud-go/gcore/securitygroup/v1/securitygroups"
	"github.com/G-Core/gcorelabscloud-go/gcore/securitygroup/v1/types"
)

func strPtr(s string) *string { return &s }

// TestFilterSecurityGroupRules_NilDescription verifies that
// filterSecurityGroupRules does not panic when a rule has a nil Description.
//
// Reproduces ICM-47063: the provider crashes during terraform plan/destroy
// when the k8s cluster's worker security group contains a rule where the API
// returns "description": null.
func TestFilterSecurityGroupRules_NilDescription(t *testing.T) {
	rules := []securitygroups.SecurityGroupRule{
		{ID: "user-rule", Direction: types.RuleDirectionIngress, Description: strPtr("allow-http")},
		{ID: "nil-desc", Direction: types.RuleDirectionEgress, Description: nil},
		{ID: "system-rule", Direction: types.RuleDirectionIngress, Description: strPtr("system")},
	}

	result := filterSecurityGroupRules(rules)

	resultIDs := make(map[string]bool, len(result))
	for _, rule := range result {
		resultIDs[rule.ID] = true
	}

	if !resultIDs["user-rule"] {
		t.Errorf("expected result to contain %q", "user-rule")
	}
	if !resultIDs["nil-desc"] {
		t.Errorf("expected result to contain %q", "nil-desc")
	}
	if resultIDs["system-rule"] {
		t.Errorf("expected result to exclude %q", "system-rule")
	}
}

// TestFilterSecurityGroupRules_FiltersSystemRules verifies that rules
// with description "system" are filtered out.
func TestFilterSecurityGroupRules_FiltersSystemRules(t *testing.T) {
	rules := []securitygroups.SecurityGroupRule{
		{ID: "user-rule", Direction: types.RuleDirectionIngress, Description: strPtr("allow-http")},
		{ID: "system-rule", Direction: types.RuleDirectionIngress, Description: strPtr("system")},
	}

	result := filterSecurityGroupRules(rules)

	if len(result) != 1 {
		t.Errorf("expected 1 rule (system filtered out), got %d", len(result))
	}
	if result[0].ID != "user-rule" {
		t.Errorf("expected remaining rule to be 'user-rule', got %q", result[0].ID)
	}
}

// TestFilterSecurityGroupRules_EmptySlice verifies no panic on empty input.
func TestFilterSecurityGroupRules_EmptySlice(t *testing.T) {
	result := filterSecurityGroupRules([]securitygroups.SecurityGroupRule{})
	if len(result) != 0 {
		t.Errorf("expected 0 rules, got %d", len(result))
	}
}

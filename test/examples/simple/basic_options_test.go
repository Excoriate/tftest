package simple

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Excoriate/tftest/pkg/scenario"
)

func TestSimpleOptionsPlanScenario(t *testing.T) {
	s, err := scenario.New(t, "../../data/tf-random")
	assert.NoErrorf(t, err, "Failed to create scenario: %s", err)

	s.Stg.PlanStage(t, s.GetTerraformOptions())
}

func TestExpectAnyChangeOnPlan(t *testing.T) {
	s, err := scenario.New(t, "../../data/tf-random")
	assert.NoErrorf(t, err, "Failed to create scenario: %s", err)

	s.Stg.PlanStageWithAnySortOfChanges(t, s.GetTerraformOptions())
}

func TestSpecificResourcesExpectedChanges(t *testing.T) {
	s, err := scenario.New(t, "../../data/tf-random")
	assert.NoErrorf(t, err, "Failed to create scenario: %s", err)

	s.Stg.PlanWithSpecificResourcesThatWillChange(t, s.GetTerraformOptions(), []string{"random_id.this"})
}

func TestLifecycle(t *testing.T) {
	s, err := scenario.New(t, "../../data/tf-random")
	assert.NoErrorf(t, err, "Failed to create scenario: %s", err)

	defer s.Stg.DestroyStage(t, s.GetTerraformOptions())

	s.Stg.PlanStageWithAnySortOfChanges(t, s.GetTerraformOptions())
	s.Stg.ApplyStage(t, s.GetTerraformOptions())
}

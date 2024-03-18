package simple

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Excoriate/tftest/pkg/scenario"
)

func TestWithVarOptionsInvalid(t *testing.T) {
	workdir := "../../data/tf-random"
	_, err := scenario.NewWithOptions(t, workdir,
		scenario.WithVarFiles(workdir, "i-do-not-exist.tfvars"))

	assert.Errorf(t, err, "It was expected to fail with an error due to the non-existent file")
	assert.ErrorContainsf(t, err, "the terraform variable file does not exist: ../../data/tf-random/i-do-not-exist.tfvars", "It was expected to fail with an error due to the non-existent file")
}

func TestWithVarOptionsValid(t *testing.T) {
	workdir := "../../data/tf-random"
	s, err := scenario.NewWithOptions(t, workdir,
		scenario.WithVarFiles(workdir, "fixtures/override-random-password.tfvars"))

	assert.NoErrorf(t, err, "Failed to create scenario: %s", err)

	s.Stg.PlanWithSpecificVariableValueToExpect(t, s.GetTerraformOptions(), "random_length_password", "25")
}

func TestWithVarsValid(t *testing.T) {
	workdir := "../../data/tf-random"
	s, err := scenario.NewWithOptions(t, workdir,
		scenario.WithVars(map[string]interface{}{"random_length_password": 10}))

	assert.NoErrorf(t, err, "Failed to create scenario: %s", err)

	s.Stg.PlanWithSpecificVariableValueToExpect(t, s.GetTerraformOptions(), "random_length_password", "10")
}

func TestWithParallelism(t *testing.T) {
	workdir := "../../data/tf-random"
	s, err := scenario.NewWithOptions(t, workdir, scenario.WithParallel())

	assert.NoErrorf(t, err, "Failed to create scenario: %s", err)

	s.Stg.PlanStage(t, s.GetTerraformOptions())
}

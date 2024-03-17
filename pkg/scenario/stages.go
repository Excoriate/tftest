package scenario

import (
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/require"
	"testing"
)

func DestroyStage(t *testing.T, options *terraform.Options) {
	out, err := terraform.DestroyE(t, options)
	require.NoErrorf(t, err, "Failed to destroy terraform: %s", out)
}

func ApplyStage(t *testing.T, options *terraform.Options) {
	out, err := terraform.InitAndApplyE(t, options)
	require.NoErrorf(t, err, "Failed to apply terraform: %s", out)
}

func PlanStage(t *testing.T, options *terraform.Options) {
	out, err := terraform.InitAndPlanE(t, options)
	require.NoErrorf(t, err, "Failed to plan terraform: %s", out)
}

func PlanStageWithExpectedChanges(t *testing.T, options *terraform.Options, expectedChanges int) {
	out, err := terraform.InitAndPlanAndShowWithStructE(t, options)
	require.NoErrorf(t, err, "Failed to plan terraform: %s", out)

	for _, change := range out.RawPlan.ResourceChanges {
		if change.Change.Actions.Create() {
			expectedChanges--
		}
	}

	require.Equalf(t, 0, expectedChanges, "Expected changes not found: %s", out)
}

func PlanStageWithDetailedExpectedChanges(t *testing.T, options *terraform.Options, expectedAdds int, expectedDeletes int, expectedUpdates int) {
	planStruct, err := terraform.InitAndPlanAndShowWithStructE(t, options)
	require.NoErrorf(t, err, "Failed to plan terraform: %+v", planStruct)

	actualAdds, actualDeletes, actualUpdates := 0, 0, 0
	for _, change := range planStruct.RawPlan.ResourceChanges {
		switch {
		case change.Change.Actions.Create():
			actualAdds++
		case change.Change.Actions.Delete():
			actualDeletes++
		case change.Change.Actions.Update():
			actualUpdates++
		}
	}

	require.Equalf(t, expectedAdds, actualAdds, "Expected and actual additions do not match")
	require.Equalf(t, expectedDeletes, actualDeletes, "Expected and actual deletions do not match")
	require.Equalf(t, expectedUpdates, actualUpdates, "Expected and actual updates do not match")
}

func PlanStageWithAnySortOfChanges(t *testing.T, options *terraform.Options) {
	out, err := terraform.InitAndPlanAndShowWithStructE(t, options)
	require.NoErrorf(t, err, "Failed to plan terraform: %s", out)
	require.NotEmptyf(t, out.RawPlan.ResourceChanges, "No changes found: %s", out)
}

func PlanStageExpectedNoChanges(t *testing.T, options *terraform.Options) {
	out, err := terraform.InitAndPlanAndShowWithStructE(t, options)
	require.NoErrorf(t, err, "Failed to plan terraform: %s", out)
	require.Emptyf(t, out.RawPlan.ResourceChanges, "Changes found: %s", out)
}

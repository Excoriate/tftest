package scenario

import (
	"fmt"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/stretchr/testify/assert"
	"reflect"
	"strconv"
	"testing"

	tfjson "github.com/hashicorp/terraform-json"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/require"
)

type StageClient struct{}

type TestType int

const (
	ShouldContain = iota
	ShouldNOTContain
	ShouldBeEqual
)

type JSONPathTestCases struct {
	TestName           string
	ExpectedValue      interface{}
	JSONPathToCompare  string
	AllowDifferentType bool
	TestType           TestType
}

type Stage interface {
	DestroyStage(t *testing.T, options *terraform.Options)
	PlanStage(t *testing.T, options *terraform.Options)
	ApplyStage(t *testing.T, options *terraform.Options)
	PlanStageWithExpectedChanges(t *testing.T, options *terraform.Options, expectedChanges int)
	PlanStageWithDetailedExpectedChanges(t *testing.T, options *terraform.Options, expectedAdds, expectedDeletes, expectedUpdates int)
	PlanStageWithAnySortOfChanges(t *testing.T, options *terraform.Options)
	PlanStageExpectedNoChanges(t *testing.T, options *terraform.Options)
	PlanWithSpecificResourcesThatWillChange(t *testing.T, options *terraform.Options, resources []string)
	PlanWithResourcesExpectedToBeCreated(t *testing.T, options *terraform.Options, resources []string)
	PlanWithResourcesExpectedToBeDeleted(t *testing.T, options *terraform.Options, resources []string)
	PlanWithResourcesExpectedToBeUpdated(t *testing.T, options *terraform.Options, resources []string)
	PlanWithSpecificVariableValueToExpect(t *testing.T, options *terraform.Options, variable, value string)
}

func (c *StageClient) DestroyStage(t *testing.T, options *terraform.Options) {
	out, err := terraform.DestroyE(t, options)
	require.NoErrorf(t, err, "Failed to destroy terraform: %s", out)
}

func (c *StageClient) PlanStage(t *testing.T, options *terraform.Options) {
	out, err := terraform.InitAndPlanE(t, options)
	require.NoErrorf(t, err, "Failed to plan terraform: %s", out)
}

func (c *StageClient) ApplyStage(t *testing.T, options *terraform.Options) {
	out, err := terraform.InitAndApplyE(t, options)
	require.NoErrorf(t, err, "Failed to apply terraform: %s", out)
}

func (c *StageClient) PlanStageWithExpectedChanges(t *testing.T, options *terraform.Options, expectedChanges int) {
	out, err := terraform.InitAndPlanAndShowWithStructE(t, options)
	require.NoErrorf(t, err, "Failed to plan terraform: %s", out)

	for _, change := range out.RawPlan.ResourceChanges {
		if change.Change.Actions.Create() {
			expectedChanges--
		}
	}

	require.Equalf(t, 0, expectedChanges, "Expected and actual changes do not match")
}

func (c *StageClient) PlanStageWithDetailedExpectedChanges(t *testing.T, options *terraform.Options, expectedAdds, expectedDeletes, expectedUpdates int) {
	planStruct, err := terraform.InitAndPlanAndShowWithStructE(t, options)
	require.NoErrorf(t, err, "Failed to plan terraform: %s", planStruct)

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

func (c *StageClient) PlanStageWithAnySortOfChanges(t *testing.T, options *terraform.Options) {
	out, err := terraform.InitAndPlanAndShowWithStructE(t, options)
	require.NoErrorf(t, err, "Failed to plan terraform: %s", out)
	require.NotEmptyf(t, out.RawPlan.ResourceChanges, "No changes found: %s", out)
	require.Truef(t, len(out.RawPlan.ResourceChanges) > 0, "No changes found: %s", out)
}

func (c *StageClient) PlanStageExpectedNoChanges(t *testing.T, options *terraform.Options) {
	out, err := terraform.InitAndPlanAndShowWithStructE(t, options)
	require.NoErrorf(t, err, "Failed to plan terraform: %s", out)
	require.Emptyf(t, out.RawPlan.ResourceChanges, "Changes found: %s", out)
}

func (c *StageClient) PlanWithSpecificResourcesThatWillChange(t *testing.T, options *terraform.Options, resources []string) {
	out, err := terraform.InitAndPlanAndShowWithStructE(t, options)
	require.NoErrorf(t, err, "Failed to plan terraform: %s", out)

	resourceChangeFound := make(map[string]bool)
	for _, resource := range resources {
		resourceChangeFound[resource] = false
	}

	for _, change := range out.RawPlan.ResourceChanges {
		if _, exists := resourceChangeFound[change.Address]; exists {
			if change.Change.Actions.Create() || change.Change.Actions.Delete() || change.Change.Actions.Update() {
				// Mark the resource as found and changed
				resourceChangeFound[change.Address] = true
			}
		}
	}

	// Verify that all specified resources are planned for change
	for resource, changed := range resourceChangeFound {
		require.Truef(t, changed, "Resource %s did not change but was expected to", resource)
	}
}

// CheckResourcesChanges checks if the specified resources have the expected changes
// The check function is used to determine if the resource has the expected change
func (c *StageClient) CheckResourcesChanges(t *testing.T, options *terraform.Options, resources []string, check func(tfjson.Actions) bool, failMsg string) {
	out, err := terraform.InitAndPlanAndShowWithStructE(t, options)
	require.NoErrorf(t, err, "Failed to plan terraform: %s", out)

	resourceFound := make(map[string]bool)
	for _, resource := range resources {
		resourceFound[resource] = false
	}

	for _, change := range out.RawPlan.ResourceChanges {
		if _, exists := resourceFound[change.Address]; exists && check(change.Change.Actions) {
			resourceFound[change.Address] = true
		}
	}

	for resource, found := range resourceFound {
		require.Truef(t, found, failMsg, resource)
	}
}

func (c *StageClient) PlanWithResourcesExpectedToBeCreated(t *testing.T, options *terraform.Options, resources []string) {
	c.CheckResourcesChanges(t, options, resources, func(action tfjson.Actions) bool {
		return action.Create()
	}, "Resource %s was not marked to be created but was expected to")
}

func (c *StageClient) PlanWithResourcesExpectedToBeDeleted(t *testing.T, options *terraform.Options, resources []string) {
	c.CheckResourcesChanges(t, options, resources, func(action tfjson.Actions) bool {
		return action.Delete()
	}, "Resource %s was not marked to be deleted but was expected to")
}

func (c *StageClient) PlanWithResourcesExpectedToBeUpdated(t *testing.T, options *terraform.Options, resources []string) {
	c.CheckResourcesChanges(t, options, resources, func(action tfjson.Actions) bool {
		return action.Update()
	}, "Resource %s was not marked to be updated but was expected to")
}

func (c *StageClient) PlanWithSpecificVariableValueToExpect(t *testing.T, options *terraform.Options, variable, expectedValue string) {
	out, err := terraform.InitAndPlanAndShowWithStructE(t, options)
	require.NoErrorf(t, err, "Failed to plan terraform: %s", out)

	variableFromPlan, found := out.RawPlan.Variables[variable]
	require.Truef(t, found, "Variable %s was not found in the plan", variable)
	require.NotNilf(t, variableFromPlan, "Variable %s was found in the plan but was nil", variable)

	actualValue := variableFromPlan.Value
	compareValues(t, actualValue, expectedValue, variable)
}

// PlanAndAssertJSONWithJSONPath performs JSON path planning and assertion in Go testing.
//
// t *testing.T: Testing object
// options *terraform.Options: Terraform options
// testCases []JSONPathTestCases: Array of JSON path test cases
func (c *StageClient) PlanAndAssertJSONWithJSONPath(t *testing.T, options *terraform.Options, testCases []JSONPathTestCases) {
	jsonPlan := terraform.InitAndPlanAndShow(t, options)

	for _, testCase := range testCases {
		t.Run(testCase.TestName, func(t *testing.T) {
			var result []interface{}
			k8s.UnmarshalJSONPath(t, []byte(jsonPlan), testCase.JSONPathToCompare, &result)
			assert.NotNil(t, result)

			// if expected is slice then it's ok to compare entire slice
			if reflect.TypeOf(testCase.ExpectedValue).Kind() == reflect.Slice {
				assert.ObjectsAreEqual(testCase.ExpectedValue, result)
			} else {
				// work on the result[0]
				t.Logf("result returned raw: %v", result)
				v := result[0]
				if !testCase.AllowDifferentType {
					assert.Equal(t, reflect.TypeOf(testCase.ExpectedValue).Kind(), reflect.TypeOf(v).Kind())
				}
				applyTestType(t, testCase.TestType, v, testCase.ExpectedValue, fmt.Sprintf("JSONPATH query: %s", testCase.JSONPathToCompare))
			}
		})
	}
}

func applyTestType(t *testing.T, testType TestType, actual, expected interface{}, additionalMessage string) {
	switch testType {
	case ShouldContain:
		msg := fmt.Sprintf("Output did not contains the expected value. Expected: %v, Actual: %v, %s", expected, actual, additionalMessage)
		assert.Contains(t, actual, expected, msg)
	case ShouldNOTContain:
		msg := fmt.Sprintf("Output is expected NOT to contains the value. Expected: %v, Actual: %v, %s", expected, actual, additionalMessage)
		assert.NotContains(t, actual, expected, msg)
	case ShouldBeEqual:
		msg := fmt.Sprintf("Output did not match with expected value. Expected: %v, Actual: %v, %s", expected, actual, additionalMessage)
		assert.EqualValuesf(t, actual, expected, msg)
	}
}

func compareValues(t *testing.T, actual interface{}, expected, variableName string) {
	actualType := reflect.TypeOf(actual)
	if actualType == nil {
		require.Failf(t, "Type assertion failed", "Unable to determine the type of the variable %s", variableName)
		return
	}

	switch actualType.Kind() {
	case reflect.String:
		require.Equalf(t, expected, actual, "Variable %s does not have the expected value", variableName)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		expectedInt, err := strconv.ParseInt(expected, 10, 64)
		require.NoErrorf(t, err, "Expected value for variable %s is not an integer: %s", variableName, expected)
		require.Equalf(t, expectedInt, reflect.ValueOf(actual).Int(), "Variable %s does not have the expected value", variableName)
	case reflect.Float32, reflect.Float64:
		expectedFloat, err := strconv.ParseFloat(expected, 64)
		require.NoErrorf(t, err, "Expected value for variable %s is not a float: %s", variableName, expected)
		require.Equalf(t, expectedFloat, reflect.ValueOf(actual).Float(), "Variable %s does not have the expected value", variableName)
	case reflect.Bool:
		expectedBool, err := strconv.ParseBool(expected)
		require.NoErrorf(t, err, "Expected value for variable %s is not a boolean: %s", variableName, expected)
		require.Equalf(t, expectedBool, actual, "Variable %s does not have the expected value", variableName)
	default:
		require.Failf(t, "Unsupported type", "Variable %s has an unsupported type: %s", variableName, actualType.Kind().String())
	}
}

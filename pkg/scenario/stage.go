package scenario

import (
	"fmt"
	"reflect"
	"strconv"
	"testing"

	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/stretchr/testify/assert"

	tfjson "github.com/hashicorp/terraform-json"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/require"
)

// StageClient represents a client for managing Terraform stages.
type StageClient struct{}

// TestType represents the type of test to be performed.
type TestType int

const (
	ShouldContain = iota
	ShouldNOTContain
	ShouldBeEqual
)

// JSONPathTestCases represents the test cases for JSON path assertions.
type JSONPathTestCases struct {
	TestName           string
	ExpectedValue      interface{}
	JSONPathToCompare  string
	AllowDifferentType bool
	TestType           TestType
}

// Stage defines an interface for managing Terraform stages.
type Stage interface {
	DestroyStage(t *testing.T, options *terraform.Options)
	PlanStage(t *testing.T, options *terraform.Options)
	ApplyStage(t *testing.T, options *terraform.Options)
	PlanStageWithExpectedChanges(t *testing.T, options *terraform.Options, expectedChanges int)
	PlanStageWithDetailedExpectedChanges(t *testing.T, options *terraform.Options, expectedAdds, expectedDeletes, expectedUpdates int)
	PlanStageWithAnySortOfChanges(t *testing.T, options *terraform.Options)
	PlanStageExpectedNoChanges(t *testing.T, options *terraform.Options)
	PlanWithSpecificResourcesThatWillChange(t *testing.T, options *terraform.Options, resources []string)
	PlanWithSpecificResourcesThatShouldNotChange(t *testing.T, options *terraform.Options, resources []string)
	PlanWithResourcesExpectedToBeCreated(t *testing.T, options *terraform.Options, resources []string)
	PlanWithResourcesExpectedToBeDeleted(t *testing.T, options *terraform.Options, resources []string)
	PlanWithResourcesExpectedToBeUpdated(t *testing.T, options *terraform.Options, resources []string)
	PlanWithSpecificVariableValueToExpect(t *testing.T, options *terraform.Options, variable, value string)
	PlanAndAssertJSONWithJSONPath(t *testing.T, options *terraform.Options, testCases []JSONPathTestCases)
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

// DestroyStage destroys the Terraform stage.
//
// Parameters:
//   - t: The testing instance.
//   - options: The Terraform options.
func (c *StageClient) DestroyStage(t *testing.T, options *terraform.Options) {
	out, err := terraform.DestroyE(t, options)
	require.NoErrorf(t, err, "Failed to destroy terraform: %s", out)
}

// PlanStage plans the Terraform stage.
//
// Parameters:
//   - t: The testing instance.
//   - options: The Terraform options.
func (c *StageClient) PlanStage(t *testing.T, options *terraform.Options) {
	out, err := terraform.InitAndPlanE(t, options)
	require.NoErrorf(t, err, "Failed to plan terraform: %s", out)
}

// ApplyStage applies the Terraform stage.
//
// Parameters:
//   - t: The testing instance.
//   - options: The Terraform options.
func (c *StageClient) ApplyStage(t *testing.T, options *terraform.Options) {
	out, err := terraform.InitAndApplyE(t, options)
	require.NoErrorf(t, err, "Failed to apply terraform: %s", out)
}

// PlanStageWithExpectedChanges plans the Terraform stage and checks for the expected number of changes.
//
// Parameters:
//   - t: The testing instance.
//   - options: The Terraform options.
//   - expectedChanges: The expected number of changes.
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

// PlanStageWithDetailedExpectedChanges plans the Terraform stage and checks for the expected number of additions, deletions, and updates.
//
// Parameters:
//   - t: The testing instance.
//   - options: The Terraform options.
//   - expectedAdds: The expected number of additions.
//   - expectedDeletes: The expected number of deletions.
//   - expectedUpdates: The expected number of updates.
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

// PlanStageWithAnySortOfChanges plans the Terraform stage and checks for any sort of changes.
//
// Parameters:
//   - t: The testing instance.
//   - options: The Terraform options.
func (c *StageClient) PlanStageWithAnySortOfChanges(t *testing.T, options *terraform.Options) {
	out, err := terraform.InitAndPlanAndShowWithStructE(t, options)
	require.NoErrorf(t, err, "Failed to plan terraform: %s", out)
	require.NotEmptyf(t, out.RawPlan.ResourceChanges, "No changes found: %s", out)
	require.Truef(t, len(out.RawPlan.ResourceChanges) > 0, "No changes found: %s", out)
}

// PlanStageExpectedNoChanges plans the Terraform stage and expects no changes.
//
// Parameters:
//   - t: The testing instance.
//   - options: The Terraform options.
func (c *StageClient) PlanStageExpectedNoChanges(t *testing.T, options *terraform.Options) {
	out, err := terraform.InitAndPlanAndShowWithStructE(t, options)
	require.NoErrorf(t, err, "Failed to plan terraform: %s", out)
	require.Emptyf(t, out.RawPlan.ResourceChanges, "Changes found: %s", out)
}

// PlanWithSpecificResourcesThatWillChange plans the Terraform stage and checks that the specified resources will change.
//
// Parameters:
//   - t: The testing instance.
//   - options: The Terraform options.
//   - resources: A list of resource addresses that are expected to change.
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
				resourceChangeFound[change.Address] = true
			}
		}
	}

	for resource, changed := range resourceChangeFound {
		require.Truef(t, changed, "Resource %s did not change but was expected to", resource)
	}
}

// PlanWithSpecificResourcesThatShouldNotChange plans the Terraform stage and checks that the specified resources should not change.
//
// Parameters:
//   - t: The testing instance.
//   - options: The Terraform options.
//   - resources: A list of resource addresses that are expected not to change.
func (c *StageClient) PlanWithSpecificResourcesThatShouldNotChange(t *testing.T, options *terraform.Options, resources []string) {
	out, err := terraform.InitAndPlanAndShowWithStructE(t, options)
	require.NoErrorf(t, err, "Failed to plan terraform: %s", out)

	resourceChangeFound := make(map[string]bool)
	for _, resource := range resources {
		resourceChangeFound[resource] = false
	}

	for _, change := range out.RawPlan.ResourceChanges {
		if _, exists := resourceChangeFound[change.Address]; exists {
			if change.Change.Actions.Create() || change.Change.Actions.Delete() || change.Change.Actions.Update() {
				resourceChangeFound[change.Address] = true
			}
		}
	}

	for resource, changed := range resourceChangeFound {
		require.Falsef(t, changed, "Resource %s changed but was expected not to", resource)
	}
}

// PlanWithResourcesExpectedToBeCreated plans the Terraform stage and checks that the specified resources are expected to be created.
//
// Parameters:
//   - t: The testing instance.
//   - options: The Terraform options.
//   - resources: A list of resource addresses that are expected to be created.
func (c *StageClient) PlanWithResourcesExpectedToBeCreated(t *testing.T, options *terraform.Options, resources []string) {
	c.CheckResourcesChanges(t, options, resources, func(action tfjson.Actions) bool {
		return action.Create()
	}, "Resource %s was not marked to be created but was expected to")
}

// PlanWithResourcesExpectedToBeDeleted plans the Terraform stage and checks that the specified resources are expected to be deleted.
//
// Parameters:
//   - t: The testing instance.
//   - options: The Terraform options.
//   - resources: A list of resource addresses that are expected to be deleted.
func (c *StageClient) PlanWithResourcesExpectedToBeDeleted(t *testing.T, options *terraform.Options, resources []string) {
	c.CheckResourcesChanges(t, options, resources, func(action tfjson.Actions) bool {
		return action.Delete()
	}, "Resource %s was not marked to be deleted but was expected to")
}

// PlanWithResourcesExpectedToBeUpdated plans the Terraform stage and checks that the specified resources are expected to be updated.
//
// Parameters:
//   - t: The testing instance.
//   - options: The Terraform options.
//   - resources: A list of resource addresses that are expected to be updated.
func (c *StageClient) PlanWithResourcesExpectedToBeUpdated(t *testing.T, options *terraform.Options, resources []string) {
	c.CheckResourcesChanges(t, options, resources, func(action tfjson.Actions) bool {
		return action.Update()
	}, "Resource %s was not marked to be updated but was expected to")
}

// PlanWithSpecificVariableValueToExpect plans the Terraform stage and checks that the specified variable has the expected value.
//
// Parameters:
//   - t: The testing instance.
//   - options: The Terraform options.
//   - variable: The name of the variable to check.
//   - expectedValue: The expected value of the variable.
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
// Parameters:
//   - t: The testing object.
//   - options: The Terraform options.
//   - testCases: An array of JSON path test cases.
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

// applyTestType applies the specified test type to validate the actual value against the expected value.
//
// Parameters:
//   - t: The testing instance.
//   - testType: The type of test to be performed.
//   - actual: The actual value from the test.
//   - expected: The expected value for the test.
//   - additionalMessage: Additional message to include in the assertion.
func applyTestType(t *testing.T, testType TestType, actual, expected interface{}, additionalMessage string) {
	switch testType {
	case ShouldContain:
		msg := fmt.Sprintf("Output did not contain the expected value. Expected: %v, Actual: %v, %s", expected, actual, additionalMessage)
		assert.Contains(t, actual, expected, msg)
	case ShouldNOTContain:
		msg := fmt.Sprintf("Output is expected NOT to contain the value. Expected: %v, Actual: %v, %s", expected, actual, additionalMessage)
		assert.NotContains(t, actual, expected, msg)
	case ShouldBeEqual:
		msg := fmt.Sprintf("Output did not match the expected value. Expected: %v, Actual: %v, %s", expected, actual, additionalMessage)
		assert.EqualValuesf(t, actual, expected, msg)
	}
}

// compareValues compares the actual value against the expected value for the specified variable.
//
// Parameters:
//   - t: The testing instance.
//   - actual: The actual value from the test plan.
//   - expected: The expected value for the variable.
//   - variableName: The name of the variable being tested.
func compareValues(t *testing.T, actual interface{}, expected, variableName string) {
	actualType := reflect.TypeOf(actual)
	if actualType == nil {
		require.Failf(t, "Type assertion failed", "Unable to determine the type of the variable %s", variableName)
		return
	}

	switch actualType.Kind() {
	case reflect.String:
		assert.Equalf(t, expected, actual, "Variable %s does not have the expected value", variableName)
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

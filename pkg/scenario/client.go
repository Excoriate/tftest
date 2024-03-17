package scenario

import (
	"github.com/gruntwork-io/terratest/modules/terraform"
	test_structure "github.com/gruntwork-io/terratest/modules/test-structure"
	"testing"
)

type Options struct {
	vars     map[string]interface{}
	varFiles []string
}

type OptFn func(*Options) error

func WithVars(vars map[string]interface{}) OptFn {
	return func(o *Options) error {
		o.vars = vars
		return nil
	}
}

func WithVarFiles(varFiles ...string) OptFn {
	return func(o *Options) error {
		o.varFiles = varFiles
		return nil
	}
}

func NewWithOptions(t *testing.T, workdir string, opts ...OptFn) (*terraform.Options, error) {
	o := &Options{}
	for _, opt := range opts {
		if err := opt(o); err != nil {
			return nil, err
		}
	}
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: workdir,
		Vars:         o.vars,
		NoColor:      true,
		PlanFilePath: "plan.out",
		VarFiles:     o.varFiles,
	})

	test_structure.SaveTerraformOptions(t, workdir, terraformOptions)

	return terraformOptions, nil
}

// New creates a new Terraform options with default retryable errors and saves it to the workdir
// This is a wrapper around terraform.WithDefaultRetryableErrors
func New(t *testing.T, workdir string) *terraform.Options {
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: workdir,
		NoColor:      true,
		PlanFilePath: "plan.out",
	})

	test_structure.SaveTerraformOptions(t, workdir, terraformOptions)

	return terraformOptions
}

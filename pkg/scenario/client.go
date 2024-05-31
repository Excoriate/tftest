package scenario

import (
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/Excoriate/tftest/pkg/cloudprovider"
	"github.com/Excoriate/tftest/pkg/tfvars"
	"github.com/Excoriate/tftest/pkg/utils"
	"github.com/Excoriate/tftest/pkg/validation"
	"github.com/gruntwork-io/terratest/modules/terraform"
)

// Options represents the configuration options for a Terraform scenario.
type Options struct {
	vars         map[string]interface{}
	varFiles     []string
	enableAWS    bool
	awsRegion    string
	isParallel   bool
	retryOptions *retryableOptions
	envVars      map[string]string
	planFile     string
}

// retryableOptions represents the retry options for Terraform operations.
type retryableOptions struct {
	retryableErrors    map[string]string
	timeBetweenRetries time.Duration
	maxRetries         int
}

// OptFn is a function type used to modify Options.
type OptFn func(*Options) error

// Client represents a Terraform client for managing Terraform operations.
type Client struct {
	t        *testing.T
	opts     *terraform.Options
	Stg      *StageClient
	awsCloud cloudprovider.AWSAdapter
}

// Config defines an interface for obtaining Terraform options and AWS configuration.
type Config interface {
	GetTerraformOptions() *terraform.Options
	GetAWS() cloudprovider.AWSAdapter
}

// GetTerraformOptions returns the Terraform options for the client.
// If the options are not set, it returns an empty Terraform options object.
//
// Returns:
//   - *terraform.Options: The Terraform options.
func (c *Client) GetTerraformOptions() *terraform.Options {
	if c.opts == nil {
		return &terraform.Options{}
	}

	return c.opts
}

// GetAWS returns the AWS Cloud Provider (Client) for the client.
//
// Returns:
//   - cloudprovider.AWSAdapter: The AWS Cloud Provider (Client).
func (c *Client) GetAWS() cloudprovider.AWSAdapter {
	return c.awsCloud
}

// WithVars sets the Terraform variables for the options.
//
// Parameters:
//   - vars: A map of Terraform variables.
//
// Returns:
//   - OptFn: A function to modify the options.
func WithVars(vars map[string]interface{}) OptFn {
	return func(o *Options) error {
		o.vars = vars
		return nil
	}
}

// WithPlanFile sets the plan file path for the options.
//
// Parameters:
//   - planFile: The path to the plan file.
//
// Returns:
//   - OptFn: A function to modify the options.
func WithPlanFile(planFile string) OptFn {
	return func(o *Options) error {
		o.planFile = planFile
		return nil
	}
}

// WithAWS enables the AWS Cloud Provider (Client) for the options and sets the AWS region.
//
// Parameters:
//   - region: The AWS region. If not set, it defaults to "us-west-2".
//
// Returns:
//   - OptFn: A function to modify the options.
func WithAWS(region string) OptFn {
	return func(o *Options) error {
		if region == "" {
			region = "us-west-2"
		}

		o.enableAWS = true
		o.awsRegion = region

		return nil
	}
}

// WithRetry sets the retry options for the Terraform operations.
//
// Parameters:
//   - retryableErrors: A map of retryable errors.
//   - timeBetweenRetries: The duration to wait between retries.
//   - maxRetries: The maximum number of retries.
//
// Returns:
//   - OptFn: A function to modify the options.
func WithRetry(retryableErrors map[string]string, timeBetweenRetries time.Duration, maxRetries int) OptFn {
	return func(o *Options) error {
		o.retryOptions = &retryableOptions{
			retryableErrors:    retryableErrors,
			timeBetweenRetries: timeBetweenRetries,
			maxRetries:         maxRetries,
		}

		return nil
	}
}

// WithParallel enables parallelism for the options.
//
// Returns:
//   - OptFn: A function to modify the options.
func WithParallel() OptFn {
	return func(o *Options) error {
		o.isParallel = true
		return nil
	}
}

// WithEnvVars sets the environment variables for the options.
//
// Parameters:
//   - envVars: A map of environment variables.
//
// Returns:
//   - OptFn: A function to modify the options.
func WithEnvVars(envVars map[string]string) OptFn {
	return func(o *Options) error {
		o.envVars = envVars
		return nil
	}
}

// WithVarFiles sets the variable files for the options.
//
// Parameters:
//   - workdir: The working directory.
//   - varFiles: A list of variable file names.
//
// Returns:
//   - OptFn: A function to modify the options.
func WithVarFiles(workdir string, varFiles ...string) OptFn {
	return func(o *Options) error {
		if err := validation.IsValidTFDir(workdir); err != nil {
			return err
		}

		for _, vf := range varFiles {
			if err := validation.IsValidTFVarFile(filepath.Join(workdir, vf)); err != nil {
				return err
			}
		}

		o.varFiles = varFiles
		return nil
	}
}

// WithScannedTFVars scans the working directory and fixtures directory for Terraform variable files
// and sets them for the options.
//
// Parameters:
//   - workdir: The working directory.
//   - fixturesDir: The fixtures directory.
//
// Returns:
//   - OptFn: A function to modify the options.
func WithScannedTFVars(workdir, fixturesDir string) OptFn {
	return func(o *Options) error {
		if err := validation.IsValidTFDir(workdir); err != nil {
			return err
		}

		fixturesDirPath := filepath.Join(workdir, fixturesDir)

		if err := validation.IsValidTFDir(fixturesDirPath); err != nil {
			return err
		}

		hasTFVars, err := validation.HasTFVarFiles(fixturesDirPath)
		if err != nil {
			return err
		}

		if !hasTFVars {
			return fmt.Errorf("the Terraform module %s with this fixtures directory %s does not have any .tfvars files", workdir, fixturesDir)
		}

		tfVarsPath, tfVarsErr := tfvars.GetTFVarsFromWorkdir(workdir)
		if tfVarsErr != nil {
			return tfVarsErr
		}

		// Add the fixtures folder on each file
		for i, tfVar := range tfVarsPath {
			tfVarsPath[i] = filepath.Join(fixturesDir, tfVar)
		}

		o.varFiles = utils.MergeSlices(o.varFiles, tfVarsPath)

		return nil
	}
}

// NewWithOptions creates a new Client with the specified options.
//
// Parameters:
//   - t: The testing instance.
//   - workdir: The working directory.
//   - opts: A list of option functions to modify the options.
//
// Returns:
//   - *Client: A new Client instance.
//   - error: An error if the Client could not be created.
func NewWithOptions(t *testing.T, workdir string, opts ...OptFn) (*Client, error) {
	o := &Options{}
	for _, opt := range opts {
		if err := opt(o); err != nil {
			return nil, err
		}
	}

	tfDir, err := GetTerraformDir(t, workdir, o.isParallel)
	if err != nil {
		return nil, err
	}

	c := &Client{}

	tfOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: tfDir,
		NoColor:      true,
	})

	if o.planFile != "" {
		tfOptions.PlanFilePath = filepath.Join(tfDir, o.planFile)
	}

	if o.enableAWS {
		cfg, err := cloudprovider.NewAWS(o.awsRegion)
		t.Logf("Enabling AWS Cloud Provider (Client) with region: %s", o.awsRegion)
		if err != nil {
			return nil, err
		}

		c.awsCloud = cfg
	}

	if len(o.vars) > 0 {
		t.Logf("Setting Terraform variables: %v", o.vars)
		tfOptions.Vars = o.vars
	}

	if len(o.varFiles) > 0 {
		t.Logf("Setting Terraform variable files: %v", o.varFiles)
		tfOptions.VarFiles = o.varFiles
	}

	if o.retryOptions != nil {
		tfOptions.RetryableTerraformErrors = o.retryOptions.retryableErrors
		tfOptions.TimeBetweenRetries = o.retryOptions.timeBetweenRetries
		tfOptions.MaxRetries = o.retryOptions.maxRetries
	}

	c.opts = tfOptions

	return c, nil
}

// New creates a new Terraform Client with default retryable errors and saves it to the workdir.
// This is a wrapper around terraform.WithDefaultRetryableErrors.
//
// Parameters:
//   - t: The testing instance.
//   - workdir: The working directory.
//
// Returns:
//   - *Client: A new Client instance.
//   - error: An error if the Client could not be created.
func New(t *testing.T, workdir string) (*Client, error) {
	if err := validation.IsValidTFModuleDir(workdir); err != nil {
		return nil, err
	}

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: workdir,
		NoColor:      true,
	})

	return &Client{
		t:    t,
		opts: terraformOptions,
		Stg:  &StageClient{},
	}, nil
}

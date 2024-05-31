package terratestopts

import (
	"github.com/Excoriate/tftest/pkg/utils"
	"github.com/gruntwork-io/terratest/modules/terraform"
)

// AddEnvVars adds the specified environment variables to the Terraform options.
//
// Parameters:
//   - options: The Terraform options to which environment variables will be added.
//   - envVars: A map of environment variables to add.
//
// Returns:
//   - *terraform.Options: The updated Terraform options with the added environment variables.
//
// Example:
//
//	options := &terraform.Options{}
//	envVars := map[string]string{"FOO": "bar"}
//	updatedOptions := AddEnvVars(options, envVars)
//	fmt.Printf("Updated Terraform options: %+v\n", updatedOptions)
func AddEnvVars(options *terraform.Options, envVars map[string]string) *terraform.Options {
	for key, value := range envVars {
		options.EnvVars[key] = value
	}

	return options
}

// AddEnvVarsFromHost adds all environment variables from the host to the Terraform options.
//
// Parameters:
//   - options: The Terraform options to which environment variables will be added.
//
// Returns:
//   - *terraform.Options: The updated Terraform options with the added environment variables from the host.
//
// Example:
//
//	options := &terraform.Options{}
//	updatedOptions := AddEnvVarsFromHost(options)
//	fmt.Printf("Updated Terraform options with host environment variables: %+v\n", updatedOptions)
func AddEnvVarsFromHost(options *terraform.Options) *terraform.Options {
	envVarsFromHost := utils.GetAllEnvVarsFromHost()
	return AddEnvVars(options, envVarsFromHost)
}

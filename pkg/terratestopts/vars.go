package terratestopts

import (
	"github.com/gruntwork-io/terratest/modules/terraform"
)

// OverrideVars overrides the existing Terraform variables in the options with the specified variables.
//
// Parameters:
//   - options: The Terraform options whose variables will be overridden.
//   - vars: A map of variables to override in the options.
//
// Returns:
//   - *terraform.Options: The updated Terraform options with the overridden variables.
//
// Example:
//
//	options := &terraform.Options{}
//	vars := map[string]interface{}{"instance_type": "t2.micro"}
//	updatedOptions := OverrideVars(options, vars)
//	fmt.Printf("Updated Terraform options: %+v\n", updatedOptions)
func OverrideVars(options *terraform.Options, vars map[string]interface{}) *terraform.Options {
	options.Vars = vars

	return options
}

// AddVars adds the specified variables to the existing Terraform variables in the options.
//
// Parameters:
//   - options: The Terraform options to which variables will be added.
//   - vars: A map of variables to add to the options.
//
// Returns:
//   - *terraform.Options: The updated Terraform options with the added variables.
//
// Example:
//
//	options := &terraform.Options{}
//	vars := map[string]interface{}{"instance_type": "t2.micro"}
//	updatedOptions := AddVars(options, vars)
//	fmt.Printf("Updated Terraform options: %+v\n", updatedOptions)
func AddVars(options *terraform.Options, vars map[string]interface{}) *terraform.Options {
	for key, value := range vars {
		options.Vars[key] = value
	}

	return options
}

// AddTFVars adds the specified Terraform variable files to the options.
//
// Parameters:
//   - options: The Terraform options to which variable files will be added.
//   - varFiles: A list of variable file paths to add.
//
// Returns:
//   - *terraform.Options: The updated Terraform options with the added variable files.
//
// Example:
//
//	options := &terraform.Options{}
//	updatedOptions := AddTFVars(options, "vars.tfvars", "prod.tfvars")
//	fmt.Printf("Updated Terraform options with variable files: %+v\n", updatedOptions)
func AddTFVars(options *terraform.Options, varFiles ...string) *terraform.Options {
	options.VarFiles = append(options.VarFiles, varFiles...)

	return options
}

// OverrideTFVars overrides the existing Terraform variable files in the options with the specified variable files.
//
// Parameters:
//   - options: The Terraform options whose variable files will be overridden.
//   - varFiles: A list of variable file paths to override.
//
// Returns:
//   - *terraform.Options: The updated Terraform options with the overridden variable files.
//
// Example:
//
//	options := &terraform.Options{}
//	updatedOptions := OverrideTFVars(options, "vars.tfvars", "prod.tfvars")
//	fmt.Printf("Updated Terraform options with overridden variable files: %+v\n", updatedOptions)
func OverrideTFVars(options *terraform.Options, varFiles ...string) *terraform.Options {
	options.VarFiles = varFiles

	return options
}

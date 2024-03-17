package terratestopts

import (
	"github.com/gruntwork-io/terratest/modules/terraform"
)

func OverrideVars(options *terraform.Options, vars map[string]interface{}) *terraform.Options {
	options.Vars = vars

	return options
}

func AddVars(options *terraform.Options, vars map[string]interface{}) *terraform.Options {
	for key, value := range vars {
		options.Vars[key] = value
	}

	return options
}

func AddTFVars(options *terraform.Options, varFiles ...string) *terraform.Options {
	options.VarFiles = append(options.VarFiles, varFiles...)

	return options
}

func OverrideTFVars(options *terraform.Options, varFiles ...string) *terraform.Options {
	options.VarFiles = varFiles

	return options
}

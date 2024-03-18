package terratestopts

import (
	"github.com/Excoriate/tftest/pkg/utils"
	"github.com/gruntwork-io/terratest/modules/terraform"
)

func AddEnvVars(options *terraform.Options, envVars map[string]string) *terraform.Options {
	for key, value := range envVars {
		options.EnvVars[key] = value
	}

	return options
}

func AddEnvVarsFromHost(options *terraform.Options) *terraform.Options {
	envVarsFromHost := utils.GetAllEnvVarsFromHost()
	return AddEnvVars(options, envVarsFromHost)
}

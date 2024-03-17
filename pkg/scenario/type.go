package scenario

import "github.com/gruntwork-io/terratest/modules/terraform"

type Scenario struct {
	Name        string
	Options     *terraform.Options
	ExpectedErr error
}

variable "is_enabled" {
  type        = bool
  description = <<EOF
  Whether this module will be created or not. It is useful, for stack-composite
modules that conditionally includes resources provided by this module..
EOF
  default     = true
}

variable "tags" {
  type        = map(string)
  description = "A map of tags to add to all resources."
  default     = {}
}

/*
-------------------------------------
Custom input variables
-------------------------------------
*/
variable "random_length_string" {
  type        = number
  description = "The length of the random string to generate."
  default     = 8
}

variable "random_special_characters" {
  type        = bool
  description = "Whether to include special characters in the random string."
  default     = true
}

variable "random_length_password" {
  type        = number
  description = "The length of the random password to generate."
  default     = 16
}

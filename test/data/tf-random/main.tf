resource "random_id" "this" {
  byte_length = 8
}

resource "random_password" "this" {
  length  = var.random_length_password
  special = var.random_special_characters
}

resource "random_string" "this" {
  length  = var.random_length_string
  special = var.random_special_characters
}

resource "random_uuid" "this" {
}

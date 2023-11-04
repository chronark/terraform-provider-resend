terraform {
  required_providers {
    resend = {
      source = "registry.terraform.io/chronark/resend"
    }
  }
}

provider "resend" {}


resource "resend_domain" "example_com" {
  name   = "example.com"
  region = "us-east-1"
}

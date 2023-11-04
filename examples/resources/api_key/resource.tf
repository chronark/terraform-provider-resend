terraform {
  required_providers {
    resend = {
      source = "registry.terraform.io/chronark/resend"
    }
  }
}

provider "resend" {}


resource "resend_api_key" "my_key" {
  name       = "Vercel-Web-App"
  permission = "full_access"
}

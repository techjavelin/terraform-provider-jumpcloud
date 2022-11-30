terraform {
  required_providers {
    jumpcloud = {
      source = "techjavelin/jumpcloud"
      version = "0.0.1"
    }
  }
}

variable "jumpcloud_api_key" {
  sensitive = true
}

provider "jumpcloud" {
  api_key = var.jumpcloud_api_key
}

resource "jumpcloud_usergroup" "test" {
    name = "terraform-provider-jumpcloud"
}
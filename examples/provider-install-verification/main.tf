terraform {
    required_providers {
      jumpcloud = {
        source = "techjavelin/jumpcloud"
      }
    }
}

variable "jumpcloud_api_key" {
    sensitive = true
}

provider "jumpcloud" {
    api_key = var.jumpcloud_api_key
}

resource "jumpcloud_ad" "test_com" {
    domain = "DC=test,DC=com"
}

resource "jumpcloud_devicegroup" "test" {
  name = "Test Group"
}
terraform {
  required_providers {
    jumpcloud = {
      source = "techjavelin/jumpcloud"
    }
  }
}

provider "jumpcloud" {
  api_key = "API_KEY_GOES_HERE"
}

terraform {
  required_providers {
    mattr = {
      version = "~> 0.0.1"
      source  = "antunovic.nz/synlestidae/mattr"
    }
  }
}

provider "mattr" {
  api_url          = var.mattr_api_url
  auth_url         = var.mattr_auth_url
  client_id        = var.mattr_client_id
  client_secret    = var.mattr_client_secret
  audience    = var.mattr_auth_audience
}

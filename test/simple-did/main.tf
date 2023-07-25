terraform {
  required_providers {
    mattr = {
      version = "~> 0.0.1"
      source  = "antunovic.nz/synlestidae/mattr"
    }
  }
}

provider "mattr" {
  api_url       = "http://127.0.0.1:8080/"
  auth_url      = "http://127.0.0.1:8080/auth"
  client_id     = "test-client-id" 
  client_secret = "test-secret-id"
  audience      = "all-ages"
}

resource "mattr_did" "did" {
  method = "key"
}

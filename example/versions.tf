terraform {
  required_providers {
    mattr = {
      version = "~> 0.0.1"
      source  = "antunovic.nz/synlestidae/mattr"
    }
  }
}

provider "mattr" {
  api_url       = var.ngrok_url /*"https://53f4-161-29-134-3.ngrok.io"*/
  access_token = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJodHRwOi8vbWF0dHIvdGVuYW50LWlkIjoiMDAwMDAwMDAtMDAwMC0wMDAwLTAwMDAtMDAwMDAwMDAwMDAwIiwiYXVkIjoiaHR0cDovL2xvY2FsaG9zdCIsInNjb3BlIjoidGVuYW50czpjcmVhdGUgdGVuYW50czpyZWFkIHRlbmFudHM6ZGVsZXRlIiwiaWF0IjoxNjkxNjEzNDA2fQ.OlxC8DMDGcEMcbnoyNDZfwBeW6Ti0d9dV9LaPcGGGp58hwd_lhDkTbCf9k_ZTcA_HA76D3jtf47Yy9OXhEpZNuhQN6IjhEq2glkxoJh03dTgMDKOTNMQ5PPrAKOpuXVaVTF893pN_Sai0oYSBdtTgLSmmP4LdDAZwFuWvXGv6Hfm97oY-l4yzcawezaqSoAX6SXPbGGRk_AAxYgtFqwqL76WnshSgwFZ__4GEpSqp8QQB26OTC3toSJ-LrBK0nZP823jZFiyruDUgf8WYsjvJifYZt0j-90tZFNb7d_DIRvkpDmc-qf1pDBCvVEw6wp4VHEQMnISp_Vw56Fa8iXVog"
  /*api_url       = var.mattr_api_url
  auth_url      = var.mattr_auth_url
  client_id     = var.mattr_client_id
  client_secret = var.mattr_client_secret
  audience      = var.mattr_auth_audience*/
}

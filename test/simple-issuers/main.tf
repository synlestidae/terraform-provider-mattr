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

resource "mattr_did" "simple_did" {
  method = "key"
}

resource "mattr_issuer" "simple_issuer" {
  issuer_did      = mattr_did.simple_did.id
  issuer_name     = "University Attendance Credential"
  issuer_logo_url = "https://example.edu/img/logo.png"
  issuer_icon_url = "https://example.edu/img/icon.png"
  description     = "This credential shows that the person has attended the mentioned course and attained the relevant awards."
  context = [
    "https://schema.org"
  ]
  type                = ["AlumniCredential"]
  proof_type          = "Ed25519Signature2018"
  background_color    = "#B00AA0"
  watermark_image_url = "https://example.edu/img/watermark.png"
  url                 = "https://example-university.au.auth0.com"
  scope = [
    "openid",
    "profile",
    "email"
  ]
  client_id                  = "vJ0SCKchr4XjC0xHNE8DkH6Pmlg2lkCN"
  client_secret              = "QNwfa4Yi4Im9zy1u_15n7SzWKt-9G5cdH0r1bONRpUPfN-UIRaaXv_90z8V6-OjH"
  token_endpoint_auth_method = "client_secret_post"
  claims_source              = "userInfo"
  static_request_parameters = {
    prompt : "login",
    max_age : 10000,
  }
  forwarded_request_parameters = [
    "login_hint"
  ]
  claim_mappings {
    json_ld_term = "alumniOf"
    oidc_claim   = "alumni_of"
  }
}

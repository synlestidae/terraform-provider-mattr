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

resource "mattr_claim_source" "test_claim_source" {
  name = "My claims from example.com"
  url = "https://example.com"
  authorization_type = "api-key"
  authorization_value = "6hrFDATxrG9w14QY9wwnmVhLE0Wg6LIvwOwUaxz761m1J"

  request_parameter {
    name = "property1"
    map_from = "claims.accountType"
    default_value = "student"
  }

  request_parameter {
    name = "property2"
    map_from = "claims.accountType"
    default_value = "student"
  }
}

resource "mattr_credential_web" "test_credential" {
  name = "Course credential"
  description = "This credential shows that the person has attended a course."
  type = "CourseCredential"
  additional_types = ["AlumniCredential", "EducationCredential"]

  contexts = [
    "https://json-ld.org/contexts/person.jsonld",
  ]

  issuer_name = "ABC University"
  issuer_logo_url = "https://example.edu/img/logo.png"
  issuer_icon_url = "https://example.edu/img/icon.png"

  proof_type = ["Ed25519Signature2018"]

  background_color = "#B00AA0"
  watermark_image_url = "https://example.edu/img/watermark.png"

  claim_mapping {
    name = "firstName"
    map_from = "claims.given_name"
    required = true
  }

  claim_mapping {
    name = "address"
    map_from = "claims.address.formatted"
  }

  claim_mapping {
    name = "picture"
    map_from = "claims.picture"
    default_value = "http://example.edu/img/placeholder.png"
  }

  claim_mapping {
    name = "badge"
    default_value = "http://example.edu/img/badge.png"
  }

  claim_mapping {
    name = "providerSubjectId"
    map_from = "authenticationProvider.subjectId"
  }

  claim_source_id = mattr_claim_source.test_claim_source.id

  persist = false
  revocable = true
  years = 0
  months = 3
  weeks = 0
  days = 0
  hours = 0
  minutes = 0
  seconds = 0
}

resource "mattr_authentication_provider" "test_authentication_provider" {
  url = "https://example-university.au.auth0.com"
  scope = [
    "email"
  ]

  client_id = "vJ0SCKchr4XjC0xHNE8DkH6Pmlg2lkCN"
  client_secret = "QNwfa4Yi4Im9zy1u_15n7SzWKt-9G5cdH0r1bONRpUPfN-UIRaaXv_90z8V6-OjH"
  token_endpoint_auth_method = "client_secret_post"

  claims_source = "idToken"

  static_request_parameters = {
    "prompt": "login"
  }
  forwarded_request_parameters = [
    "login_hint"
  ]

  claims_to_sync = [
    "first_name",
    "last_name",
    "email"
  ]
}

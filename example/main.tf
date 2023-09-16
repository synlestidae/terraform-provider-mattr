resource "mattr_claim_source" "test_claim_source" {
  name = "Testo claim source"
  url = "https://webhook.site/09c0daf8-8306-413a-a30f-9e52de2db41a"
  authorization_type = "api-key"
  authorization_value = "DRPWCO7eVf81mi4mGqBV-evPkIBh4s8A"

  request_parameter {
    name = "profile"
    map_from = "credentialConfiguration.profile"
  }

  request_parameter {
    name = "profile"
    map_from = "credentialConfiguration.id"
  }
}

resource "mattr_credential_web" "test_cred_config" {
  name        = "My Test Credential"
  description = "This credential for testing only."
  type        = "ThingyCredential"
  additional_types = [
  ]
  contexts = [
    "https://json-ld.org/contexts/person.jsonld",
    "https://schema.org"
  ]
  issuer_name     = "Test Issuer"
  issuer_logo_url = "https://example.edu/img/logo.png"
  issuer_icon_url = "https://example.edu/img/icon.png"

  background_color    = "#B00AA0"
  watermark_image_url = "https://example.edu/img/watermark.png"

  claim_mapping {
    name     = "firstName"
    map_from = "claims.given_name"
    required = true
  }

  claim_mapping {
    name     = "familyName"
    map_from = "claims.family_name"
    required = true
  }

  persist         = true
  revocable       = true
  claim_source_id = mattr_claim_source.test_claim_source.id
  years           = 1
  months          = 0
  weeks           = 0
  days            = 0
  hours           = 0
  minutes         = 0
  seconds         = 0
}

resource "mattr_authentication_provider" "test_provider" {
  url = "https://dev-s8my1837fpnxqetu.us.auth0.com"
  scope = ["openid"]
  client_id = "OQGBbQP5GB9YopSUOru4dBNJcVBvkzRA"
  client_secret = "1Fc0PoLDwZBw_7TwHKAoNB-60mjfNlU9JjhmrAsJO9Os_ua-8qioQLBNM8peOgrm"
}

resource "mattr_custom_domain" "test_domain" {
  name = "My NGrok"
  domain = "53f4-161-29-134-3.ngrok.io"
  logo_url = "https://mattr.global/favicon.ico"
  homepage = "https://mattr.global"
}

resource "mattr_credential_offer" "test_offer" {
  credentials = [mattr_credential_web.test_cred_config.id]
}

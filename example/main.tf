resource "mattr_did" "did" {
  method = "web"
  url    = "www.antunovic.nz"
}

resource "mattr_did" "test_did" {
  method = "web"
  url    = "test.com"
}

resource "mattr_webhook" "issue_webhook" {
  url      = "https://webhook.site/402ec72e-097e-4833-a6c4-a6ce50d8bed6"
  events   = ["OidcIssuerCredentialIssued"]
  disabled = true
}

resource "mattr_issuer" "antunovic_issuer" {
  issuer_did      = mattr_did.did.id
  issuer_logo_url = "https://example.edu/img/logo.png"
  issuer_icon_url = "https://example.edu/img/icon.png"
  name            = "University Attendance Credential"
  description     = "This credential shows that the person has attended the mentioned course and attained the relevant awards."
  context         = ["https://schema.org"]
  type            = ["AlumniCredential"]
  federated_provider = {
    url                        = "https://accounts.google.com/"
    scope                      = "openid"
    client_id                  = "UKNVhhnFUK2T0bR05R5IRLSImEw8mLCh"
    client_secret              = "LvBLr8yeVP9i8wCUY25720XNJ63zvBP-MtMSVQFiEhsFqP5uM4ORp51Owp6Vud7_"
    token_endpoint_auth_method = "client_secret_post"
    claims_source              = "userInfo"
  }
  credential_branding = {
    background_color    = "#B00AA0"
    watermark_image_url = "https://example.edu/img/watermark.png"
  }
  static_request_parameters = {
    prompt  = "login"
    max_age = "10000"
  }
  forwarded_request_parameters = ["login_hint"]

  claim_mappings {
    json_ld_term = "alumniOf"
    oidc_claim   = "alumni_of"
  }
}

resource "mattr_credential" "antunovic_credential" {
  name = "Course credential"
  description = "This credential shows that the person has attended a course."
  type = "CourseCredential"
  additional_types = [
    "AlumniCredential",
    "EducationCredential"
  ]
  contexts = [
    "https://json-ld.org/contexts/person.jsonld"
  ]
  issuer = {
    "name": "ABC University",
    "logo_url": "https://example.edu/img/logo.png",
    "icon_url": "https://example.edu/img/icon.png"
  }
  credential_branding = {
    "background_color": "#B00AA0",
    "watermark_image_url": "https://example.edu/img/watermark.png"
  }
  claim_mappings = {
    first_name = {
      map_from = "claims.given_name"
      required = true
    }
    address = {
      map_from = "claims.address.formatted"
    }
    picture = {
      map_from = "claims.picture"
      default_value = "http://example.edu/img/placeholder.png"
    }
    badge = {
      default_value = "http://example.edu/img/badge.png"
    }
    provider_subject_id = {
      map_from = "authenticationProvider.subjectId"
    }
  }
  persist = false
  revocable = true
  claim_source_id = "78e1b90c-401d-45bb-89c0-938da4d44c60"
  expires_in = {
    years = 0
    months = 3
    weeks = 0
    days = 0
    hours = 0
    minutes = 0
    seconds = 0
  }
}

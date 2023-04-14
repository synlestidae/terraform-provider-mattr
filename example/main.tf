resource "mattr_did" "did" {
  method = "web"
  url    = "antunovic.nz"
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
    url                        = "https://accounts.google.com"
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
    max_age = 10000
  }
  forwarded_request_parameters = ["login_hint"]

  claim_mappings {
    json_ld_term = "alumniOf"
    oidc_claim   = "alumni_of"
  }
}


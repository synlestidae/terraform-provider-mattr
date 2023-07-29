/*resource "mattr_did" "did" {
  method = "key"
}

resource "mattr_claim_source" "test_claim_source" {
  name = "My claims from testo-claimo.claims.studio"
  url  = "https://testo-claimo.claims.studio/api/data"

  authorization_type  = "api-key"
  authorization_value = "your-api-key-here"

  request_parameter {
    property      = "given_name"
    map_from      = "claims.given_name"
    default_value = "Firsty"
  }

  request_parameter {
    property      = "family_name"
    map_from      = "claims.family_name"
    default_value = "Lastyname"
  }
}

resource "mattr_credential" "test_credential" {
  name        = "Test credential"
  description = "This credential shows information about a person."
  type        = "PersonCredential"
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

  proof_type = "Ed25519Signature2018"
}

resource "mattr_credential_offer" "test_credential_offer" {
  credentials = [
    mattr_credential.test_credential.id
  ]
}

resource "mattr_custom_domain" "site_worker_custom_domain" {
  name     = "Site Worker Cert"
  logo_url = "https://s3.siteworkercert.com/logo2.jpg"
  domain   = "certs.siteworkercert.com"
  homepage = "https://siteworkercert.com"
}

resource "mattr_issuer" "test_issuer" {
  name            = "Qualification Credential"
  issuer_name     = "Example University"
  issuer_did      = "did:key:z6MkjBWPPa1njEKygyr3LR3pRKkqv714vyTkfnUdP6ToFSH5"
  issuer_logo_url = "https://example.edu/img/logo.png"
  issuer_icon_url = "https://example.edu/img/icon.png"
  description     = "This credential shows that the person has attended the mentioned course and attained the relevant awards."
  context         = ["https://schema.org"]
  type            = ["AlumniCredential"]
  proof_type      = "Ed25519Signature2018"

  background_color    = "#B00AA0"
  watermark_image_url = "https://example.edu/img/watermark.png"

  url                        = "https://example-university.au.auth0.com"
  scope                      = ["openid", "profile", "email"]
  client_id                  = "vJ0SCKchr4XjC0xHNE8DkH6Pmlg2lkCN"
  client_secret              = "QNwfa4Yi4Im9zy1u_15n7SzWKt-9G5cdH0r1bONRpUPfN-UIRaaXv_90z8V6-OjH"
  token_endpoint_auth_method = "client_secret_post"
  claims_source              = "userInfo"

  static_request_parameters = {
    prompt  = "login"
    max_age = 10000
  }

  forwarded_request_parameters = ["login_hint"]

  claim_mappings {
    json_ld_term = "alumniOf"
    oidc_claim   = "alumni_of"
  }
}*/

resource "mattr_compact_credential_template" "compact_credential_template" {
  name = "Test Compact Credential"
  template_path = "template.pdf"
  file_name = "certificate.pdf"
  metadata = {
    title = "Certificate"
  }

  fonts {
    name = "PublicSans-Bold"
    file_name = "fonts/PublicSans-Bold.ttf"
  }

  fonts {
    name = "PublicSans-Regular"
    file_name = "fonts/PublicSans-Regular.ttf"
  }

  fields {
    key = "name"
    value = "{{name}}"
    is_required = true
    alternative_text = "NAME" 
    font_name = "PublicSans-Bold"
  }
}


resource "mattr_semantic_compact_credential_template" "semantic_compact_credential_template" {
  name = "Test Compact Credential"
  template_path = "template.pdf"
  file_name = "certificate.pdf"

  metadata = {
    title = "Certificate of Completion" 
  }

  fonts {
    name = "PublicSans-Bold"
    file_name = "fonts/PublicSans-Bold.ttf"
  }

  fonts {
    name = "PublicSans-Regular"
    file_name = "fonts/PublicSans-Regular.ttf"
  }

  fields {
    key = "name"
    value = "{{name}}"
    is_required = true
    alternative_text = "NAME" 
    font_name = "PublicSans-Bold"
  }
}

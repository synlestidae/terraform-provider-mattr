resource "mattr_did" "did" {
  method = "key"
  //key_type = "Bls12381G2"
}

resource "mattr_webhook" "issue_webhook" {
  url      = "https://webhook.site/402ec72e-097e-4833-a6c4-a6ce50d8bed6"
  events   = ["OidcIssuerCredentialIssued"]
  disabled = true
}


resource "mattr_compact_credential_template" "test_cc_template" {
  file_name     = "{{ code }}_{{ name }}"
  template_path = "template.pdf"
  name =	"SamplePDF_WorkingAtHightsVC"
  font_paths = [
    "fonts/PublicSans-Bold.ttf",
    "fonts/PublicSans-Regular.ttf"
  ]

  metadata = {
    title = "{{ certificationName }} Certification – {{ name }}"
  }

  fonts {
    name      = "PublicSans-Regular"
    file_name = "PublicSans-Regular.ttf"
  }

  fonts {
    name      = "PublicSans-Bold"
    file_name = "PublicSans-Bold.ttf"
  }

  fields {
    key               = "name"
    value             = "{{ name }}"
    is_required       = true
    alternative_text = "{{ name }}"
    font_name         = "PublicSans-Regular"
  }

  fields {
    key               = "code"
    value             = "{{ code }}"
    is_required       = true
    alternative_text = "{{ code }}"
    font_name         = "PublicSans-Bold"
  }

  fields {
    key               = "certificationName"
    value             = "{{ certificationName }}"
    is_required       = true
    alternative_text = "{{ certificationName }}"
    font_name         = "PublicSans-Bold"
  }

  fields {
    key               = "certificationLevel"
    value             = "{{ certificationLevel }}"
    is_required       = true
    alternative_text = "{{ certificationLevel }}"
    font_name         = "PublicSans-Regular"
  }

  fields {
    key               = "expiry"
    value             = "{{ date expiry 'dd MMM yyyy' }}"
    is_required       = true
    alternative_text = "{{ date expiry 'PPP' }}"
    font_name         = "PublicSans-Regular"
  }
}

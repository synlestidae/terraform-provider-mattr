resource "mattr_did" "did" {
  method = "web"
  url    = "antunovic.nz"
}

resource "mattr_webhook" "issue_webhook" {
  url    = "https://webhook.site/402ec72e-097e-4833-a6c4-a6ce50d8bed6"
  events = ["OidcIssuerCredentialIssued"]
  disabled = true
}

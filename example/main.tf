resource "mattr_did" "did" {
  method = "key"
  //key_type = "Bls12381G2"
}

resource "mattr_webhook" "issue_webhook" {
  url      = "https://webhook.site/402ec72e-097e-4833-a6c4-a6ce50d8bed6"
  events   = ["OidcIssuerCredentialIssued"]
  disabled = true
}

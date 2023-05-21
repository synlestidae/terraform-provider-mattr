# A third-party Terraform provider for MATTR

[MATTR VII](https://mattr.global/platform/vii/) is a platform for developing applications to use decentralised identity.
MATTR provides a [REST API](https://learn.mattr.global/api-reference) which lets you manage resources such as DIDs, 
issuers, document signers, and so on. I made this Terraform plugin to make it easier to manage your MATTR set up, and 
save time. This provider allows you to declare and manage resources without having to use an HTTP client.

# Feature progress

| Feature                   | Description                                             | Implementation   | Unit tests    | Integration tests |
| ------------------------- | --------------------------------------------------------| ---------------- | ------------- | ----------------- |
| DID resource              | Creating DIDs for issuing                               | `████████ 100 %` | `░░░░░░░░ 0%` | `░░░░░░░░ 0%`     |
| Webhook resource          | Webhooks to receive events                              | `████████ 100 %` | `░░░░░░░░ 0%` | `░░░░░░░░ 0%`     |
| Credential config resource| Specifies a credential to issue                         | `██████░░ 75 %`  | `░░░░░░░░ 0%` | `░░░░░░░░ 0%`     |
| Claim source config       | Used by credential config to retrieve claims            | `████████ 100 %` | `░░░░░░░░ 0%` | `░░░░░░░░ 0%`     |
| Authentication provider   | Just an authorization server for credential config      | `█████░░░ 75 %`  | `░░░░░░░░ 0%` | `░░░░░░░░ 0%`     |
| Issuer resource           | Specify credentials to issue                            | `██████░░ 75 %`  | `░░░░░░░░ 0%` | `░░░░░░░░ 0%`     |
| Issuer clients resource   | Configure credentials to issue in wallet                | `█░░░░░░░ 12.5 %`| `░░░░░░░░ 0%` | `░░░░░░░░ 0%`     |
| Verifier resource         | Specify which credentials are valid                     | `░░░░░░░░ 0 %`   | `░░░░░░░░ 0%` | `░░░░░░░░ 0%`     |
| Verifier client resource  | Configure presentations for valid creds                 | `░░░░░░░░ 0 %`   | `░░░░░░░░ 0%` | `░░░░░░░░ 0%`     |
| Custom domain resource    | Configure issuer and verifier pages under your domain   | `░░░░░░░░ 0 %`   | `░░░░░░░░ 0%` | `░░░░░░░░ 0%`     |

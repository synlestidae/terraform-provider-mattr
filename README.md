# A third-party Terraform provider for MATTR

[MATTR VII](https://mattr.global/platform/vii/) is a platform for developing applications to use decentralised identity.
MATTR provides a [REST API](https://learn.mattr.global/api-reference) which lets you manage resources such as DIDs, 
issuers, document signers, and so on. I made this Terraform plugin to make it easier to manage your MATTR set up, and 
save time. This provider allows you to declare and manage resources without having to use an HTTP client.

This project is independent, fan-made software intended to be used with [MATTR](https://mattr.global)'s platform. It
is not officially endorsed or sponsored by MATTR.

# Feature progress

| Feature                   | Description                                             | Implementation   | Unit tests    | Integration tests |
| ------------------------- | --------------------------------------------------------| ---------------- | ------------- | ----------------- |
| DID                       | Creating DIDs for issuing                               | `████████ 100 %` | `░░░░░░░░ 0%` | `░░░░░░░░ 0%`     |
| Webhook                   | Webhooks to receive events                              | `████████ 100 %` | `░░░░░░░░ 0%` | `░░░░░░░░ 0%`     |
| Credential config         | Specifies a credential to issue                         | `██████░░ 75 %`  | `░░░░░░░░ 0%` | `░░░░░░░░ 0%`     |
| Claim source
| Claim source config       | Used by credential config to retrieve claims            | `████████ 100 %` | `░░░░░░░░ 0%` | `░░░░░░░░ 0%`     |
| Authentication provider   | Just an authorization server for credential config      | `█████░░░ 75 %`  | `░░░░░░░░ 0%` | `░░░░░░░░ 0%`     |
| Issuer resource           | Specify credentials to issue                            | `██████░░ 75 %`  | `░░░░░░░░ 0%` | `░░░░░░░░ 0%`     |
| Issuer clients            | Configure credentials to issue in wallet                | `██████░░ 75 %`  | `░░░░░░░░ 0%` | `░░░░░░░░ 0%`     |
| Verifier resource         | Specify which credentials are valid                     | `███░░░░░ 37.5 %`| `░░░░░░░░ 0%` | `░░░░░░░░ 0%`     |
| Verifier client           | Configure presentations for valid creds                 | `███░░░░░ 37.5 %`| `░░░░░░░░ 0%` | `░░░░░░░░ 0%`     |
| Custom domain resource    | Configure issuer and verifier pages under your domain   | `█░░░░░░░ 12.5 %`| `░░░░░░░░ 0%` | `░░░░░░░░ 0%`     |
| PDF template              | PDF template based on a .zip file                       | `░░░░░░░░ 0 %`   | `░░░░░░░░ 0 %`| `░░░░░░░░ 0%`     |
| Digital pass template     | Template for Apple Pass                                 | `░░░░░░░░ 0 %`   | `░░░░░░░░ 0 %`| `░░░░░░░░ 0%`     |
| Apple Pass template       | Template for Apple Pass                                 | `░░░░░░░░ 0 %`   | `░░░░░░░░ 0 %`| `░░░░░░░░ 0%`     |
| Compact credential PDF    | Template for compact credential PDF                     | `█░░░░░░░ 12.5`  | `░░░░░░░░ 0 %`| `░░░░░░░░ 0%`     |

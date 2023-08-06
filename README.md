# A third-party Terraform provider for MATTR

[MATTR VII](https://mattr.global/platform/vii/) is a platform for developing applications to use decentralised identity.
MATTR provides a [REST API](https://learn.mattr.global/api-reference) which lets you manage resources such as DIDs, 
issuers, document signers, and so on. I made this Terraform plugin to make it easier to manage your MATTR set up, and 
save time. This provider allows you to declare and manage resources without having to use an HTTP client.

This project is independent, fan-made software intended to be used with [MATTR](https://mattr.global)'s platform. It
is not officially endorsed or sponsored by MATTR.

# Installing

To build and install the project locally:

1. Make sure you have Go installed. Terraform providers are written in Go. brew install go should do it
2. In the root of the repo, run ./build.sh.
3. Then run ./deploy.sh. This will just copy it to a location on your machine that Terraform can find

# Try it out

Check out the [example project](./example/).

# Set-up

For the provider to work, it will need to authenticate with your secret MATTR credentials. You will have received some credentials like these:

```json
{
    "tenantSubdomain": "https://YOUR_TENANT_DOMAIN.vii.mattr.global",
    "tenantId": "d97H2lc0-3ht5-4bea-ht92-41052f047440",
    "url": "https://auth.mattr.global/oauth/token",
    "audience": "https://vii.mattr.global",
    "client_id": "YOUR_CLIENT_ID",
    "client_secret": "YOUR_CLIENT_SECRET"
}
```

You define them in your Terraform provider block:

```terraform
provider "mattr" {
  api_url       = "https://YOUR_TENANT_DOMAIN.vii.mattr.global"
  auth_url      = "https://auth.mattr.global/oauth/token"
  client_id     = "YOUR_CLIENT_ID"
  client_secret = "YOUR_CLIENT_SECRET"
  audience      = "https://vii.mattr.global"
}
```

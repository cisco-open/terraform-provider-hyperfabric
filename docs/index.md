---
subcategory: ""
layout: "hyperfabric"
page_title: "Provider: Cisco Nexus Hyperfabric"
sidebar_current: "docs-hyperfabric-index"
description: |-
  The  Cisco Nexus Hyperfabric provider is used to interact with the resources provided by Cisco Nexus Hyperfabric.
  The provider needs to be configured with the proper credentials before it can be used.
---

# Cisco Nexus Hyperfabric
Cisco Nexus Hyperfabric allows you to design, deploy, and operate your Data Center network fabrics from a cloud-based service and exposes an extensive REST API that allows integrations such as the use of this Terraform provider.

# The Cisco Nexus Hyperfabric Provider
The Cisco Nexus Hyperfabric Terraform provider is used to interact with resources provided by Cisco Nexus Hyperfabric. The provider needs to be configured with proper credentials to authenticate with Cisco Nexus Hyperfabric.

## Authentication
The Cisco Nexus Hyperfabric API requires access via an authenticated and authorized account. Only authorized accounts are able to submit requests to API operations. All operations must communicate over a secure HTTPS connection.

Before you can use the Cisco Nexus Hyperfabric API, you must first log in to the Cisco Nexus Hyperfabric service using a browser and your Cisco Connection Online (CCO) credentials. In the Cisco Nexus Hyperfabric service, you will generate a [bearer token](https://www.rfc-editor.org/rfc/rfc6750), as described in [Getting Started](https://devnetapps.cisco.com/docs/hyperfabric-api-documentation/getting-started).

The generated bearer token authenticates the account it was created with, and only for operations within the organization in which the account was logged in when it created the token. If the account is a member of multiple organizations, you must select a specific organization and create a token for that organization's API. The platform also enforces token-specific authorization and privileges based on the bearer token's scope. To use the resources in this provider, the provided bearer token should have a scope of `READ_WRITE` or `ADMIN`.



## Example Usage

```hcl
terraform {
  required_providers {
    hyperfabric = {
      source = "cisco-open/hyperfabric"
    }
  }
}

provider "hyperfabric" {
  # Configuration options
  ## Use the HYPERFABRIC_TOKEN environment variable to set your bearer token or use:
  ## token = "YOUR_HYPERFABRIC_BEARER_TOKEN"
}

resource "hyperfabric_fabric" "example" {
  name = "example-fabric"
}
```

## Schema

### Required
- `token` (string) A bearer token from the account to use for authenticating to the Cisco Nexus Hyperfabric service. See the [Getting Started](https://devnetapps.cisco.com/docs/hyperfabric-api-documentation/getting-started) page on Cisco DevNet for more information.
  - Environment variable: `HYPERFABRIC_TOKEN`

### Optional

- `proxy_creds` - (string) Proxy server credentials in the form of username:password.
  - Environment variable: `HYPERFABRIC_PROXY_CREDS`
- `proxy_url` - (string) Proxy Server URL with port number.
  - Environment variable: `HYPERFABRIC_PROXY_URL`
- `retries` - (integer) Number of retries for REST API calls.
  - Default: `2`
  - Environment variable: `HYPERFABRIC_RETRIES`
- `label` - (string) Global label for the provider.
  - Default: `terraform`
  - Environment variable: `HYPERFABRIC_LABEL`
- `url` - (string) URL of the Cisco Nexus Hyperfabric service.
  - Default: `https://hyperfabric.cisco.com`
  - Environment variable: `HYPERFABRIC_URL`
- `insecure` - (bool) Allow insecure HTTPS client.
  - Default: `false`
  - Environment variable: `HYPERFABRIC_INSECURE`

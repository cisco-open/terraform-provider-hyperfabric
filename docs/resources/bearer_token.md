---
subcategory: "Administration"
layout: "hyperfabric"
page_title: "Nexus Hyperfabric: hyperfabric_bearer_token"
sidebar_current: "docs-hyperfabric-resource-hyperfabric_bearer_token"
description: |-
  Manages a Nexus Hyperfabric Bearer Token
---

# hyperfabric_bearer_token

Manages a Nexus Hyperfabric Bearer Token

A Bearer Token is a JSON Web Token (JWT) used for authentication and authorization against the Cisco Nexus Hyperfabric REST API. The JWT is a compact, URL-safe means of representing a JSON object containing a set of key-value pairs as described in RFC 7159. It is passed as Bearer token in the Authentication header of every HTTP API request.

## API Paths ##

* `/bearerTokens` `POST`
* `/bearerTokens/{tokenId|name}` `GET, PUT, DELETE`

## GUI Information ##

* Location: `> {user_email} (Top Right) > API bearer tokens`

## Example Usage ##

The configuration snippet below creates a Bearer Token with only the required attributes.

```hcl
resource "hyperfabric_bearer_token" "example_bearer_token" {
  name = "my-example-token"
}
```
The configuration snippet below shows all possible attributes of a Bearer Token.

```hcl
resource "hyperfabric_bearer_token" "full_example_bearer_token" {
  name        = "my-full-example-token"
  description = "This is a Cisco Nexus Hyperfabric Bearer Token"
  not_after   = "2025-09-03T08:00:00.000Z"
  not_before  = "2024-09-03T08:00:00.000Z"
  scope       = "ADMIN"
}
```

## Schema ##

### Required ###

* `name` - (string) The name of the Bearer Token.

### Optional ###
  
* `description` - (string) The description is a user defined field to store notes about the Bearer Token.
* `not_after` - (string) The end date for the validity of the Bearer Token in [RFC3339](https://datatracker.ietf.org/doc/html/rfc3339#section-5.8) format (e.g. `YYYY-MM-DDTHH:MM:SSZ`).
  - Default:  Current time.
* `not_before` - (string) The start date for the validity of the Bearer Token in [RFC3339](https://datatracker.ietf.org/doc/html/rfc3339#section-5.8) format (e.g. `YYYY-MM-DDTHH:MM:SSZ`).
  - Default:  30 days from `not_before` date.
* `scope` - (string) The scope defines the level of privilege assigned to the Bearer Token.
  - Default: `READ_ONLY`
  - Valid Values: `ADMIN`, `READ_WRITE`, `READ_ONLY`.
<!-- * `labels` - (list of strings) A list of user-defined labels that can be used for grouping and filtering objects. -->
<!-- * `annotations` - (list of maps) A list of key-value annotations to store user-defined data including complex data such as JSON.

  #### Required ####

  * `name` - (string) The name used to uniquely identify the annotation.
  * `value` - (string) The value of the annotation.

  #### Optional ####

  * `data_type` - (string) The type of data stored in the value of the annotation.
      - Default: `STRING`
      - Valid Values: `STRING`, `INT32`, `UINT32`, `INT64`, `UINT64`, `BOOL`, `TIME`, `UUID`, `DURATION`, `JSON`. -->

### Read-Only ###

* `id` - (string) The unique identifier (id) of the Bearer Token.
* `token_id` - (string) The unique identifier (id) of the Bearer Token.
* `token` - (sensitive, string) The JWT token that represent the Bearer Token.
* `metadata` - (map) A map of the Metadata of the Bearer Token:
  * `created_at` - (string) The timestamp when this object was created in [RFC3339](https://datatracker.ietf.org/doc/html/rfc3339#section-5.8) format.
  * `created_by` - (string) The user that created this object.
  * `modified_at` - (string) The timestamp when this object was last modified in [RFC3339](https://datatracker.ietf.org/doc/html/rfc3339#section-5.8) format.
  * `modified_by` - (string) The user that modified this object last.
  * `revision_id` - (string) An integer that represent the current revision of the object.

## Importing

An existing Bearer Token can be [imported](https://www.terraform.io/docs/import/index.html) into this resource using the following command:

```bash
terraform import hyperfabric_bearer_token.example_bearer_token {tokenId|name}
```

Starting in Terraform version 1.5, an existing Bearer Token can be imported
using [import blocks](https://developer.hashicorp.com/terraform/language/import) via the following configuration:

```hcl
import {
  id = "{tokenId|name}"
  to = hyperfabric_bearer_token.example_bearer_token
}
```

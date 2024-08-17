---
subcategory: "Administration"
layout: "hyperfabric"
page_title: "Nexus Hyperfabric: hyperfabric_bearer_token"
sidebar_current: "docs-hyperfabric-data-source-hyperfabric_bearer_token"
description: |-
  Data source for a Nexus Hyperfabric Bearer Token
---

# hyperfabric_bearer_token

Data source for a Nexus Hyperfabric Bearer Token

A Bearer Token is a JSON Web Token (JWT) used for authentication and authorization against the Cisco Nexus Hyperfabric REST API. The JWT is a compact, URL-safe means of representing a JSON object containing a set of key-value pairs as described in RFC 7159. It is passed as Bearer token in the Authentication header of every HTTP API request.

## API Paths ##

* `/bearerTokens/{tokenId|name}` `GET`

## GUI Information ##

* Location: `> {user_email} (Top Right) > API bearer tokens`

## Example Usage ##

```hcl
data "hyperfabric_bearer_token" "example_bearer_token" {
  name = "my-example-token"
}
```

## Schema ##

### Required ###

* `name` - (string) The name of the Bearer Token.

### Read-Only ###

* `id` - (string) The unique identifier (id) of the Bearer Token.
* `description` - (string) The description is a user defined field to store notes about the Bearer Token.
* `token_id` - (string) The unique identifier (id) of the Bearer Token.
* `token` - (sensitive, string) The JWT token that represent the Bearer Token.
* `metadata` - (map) A map of the Metadata of the Bearer Token:
  * `created_at` - (string) The timestamp when this object was created in [RFC3339](https://datatracker.ietf.org/doc/html/rfc3339#section-5.8) format.
  * `created_by` - (string) The user that created this object.
  * `modified_at` - (string) The timestamp when this object was last modified in [RFC3339](https://datatracker.ietf.org/doc/html/rfc3339#section-5.8) format.
  * `modified_by` - (string) The user that modified this object last.
  * `revision_id` - (string) An integer that represent the current revision of the object.
* `not_after` - (string) The end date for the validity of the Bearer Token in [RFC3339](https://datatracker.ietf.org/doc/html/rfc3339#section-5.8) format.
* `not_before` - (string) The start date for the validity of the Bearer Token in [RFC3339](https://datatracker.ietf.org/doc/html/rfc3339#section-5.8) format.
* `scope` - (string) The scope defines the level of privilege assigned to the Bearer Token.
  - Default: `READ_ONLY`
  - Possible Values: `ADMIN`, `READ_WRITE`, `READ_ONLY`.
<!-- * `labels` - (list of strings) A list of user-defined labels that can be used for grouping and filtering objects. -->
<!-- * `annotations` - (list of maps) A list of key-value annotations to store user-defined data including complex data such as JSON.

  #### Required ####

  * `name` - (string) The name used to uniquely identify the annotation.
  * `value` - (string) The value of the annotation.

  #### Optional ####

  * `data_type` - (string) The type of data stored in the value of the annotation.
      - Default: `STRING`
      - Valid Values: `STRING`, `INT32`, `UINT32`, `INT64`, `UINT64`, `BOOL`, `TIME`, `UUID`, `DURATION`, `JSON`. -->

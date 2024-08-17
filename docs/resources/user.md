---
subcategory: "Administration"
layout: "hyperfabric"
page_title: "Nexus Hyperfabric: hyperfabric_user"
sidebar_current: "docs-hyperfabric-resource-hyperfabric_user"
description: |-
  Manages a Nexus Hyperfabric User
---

# hyperfabric_user

Manages a Nexus Hyperfabric User

A User is a Cisco.com account authorized to access a specific Cisco Nexus Hyperfabric Organization with a specific role that represents the level of privilege of the User.

## API Paths ##

* `/users` `POST`
* `/users/{userId|email}` `GET, PUT, DELETE`

## GUI Information ##

* Location: `> Administration > User management`

## Example Usage ##

The configuration snippet below creates a User with only the required attributes.

```hcl
resource "hyperfabric_user" "example_user" {
  email = "my-example-user@mydomain.mytld"
}
```
The configuration snippet below shows all possible attributes of a User.

```hcl
resource "hyperfabric_user" "full_example_user" {
  email   = "my-full-example-user@mydomain.mytld"
  enabled = true
  role    = "ADMIN"
  labels  = [
    "sj01-1-101-AAA01",
    "blue"
  ]
}
```

## Schema ##

### Required ###

* `email` - (string) The email of the User.

### Optional ###

* `enabled` - (bool) The enabled state of the User.
  - Default: `true`
* `role` - (string) The role assigned to the User that represents the level of privilege of the User.
  - Default: `READ_ONLY`
  - Valid Values: `ADMIN`, `READ_WRITE`, `READ_ONLY`.
* `labels` - (list of strings) A list of user-defined labels that can be used for grouping and filtering objects.
<!-- * `annotations` - (list of maps) A list of key-value annotations to store user-defined data including complex data such as JSON.

  #### Required ####

  * `name` - (string) The name used to uniquely identify the annotation.
  * `value` - (string) The value of the annotation.

  #### Optional ####

  * `data_type` - (string) The type of data stored in the value of the annotation.
      - Default: `STRING`
      - Valid Values: `STRING`, `INT32`, `UINT32`, `INT64`, `UINT64`, `BOOL`, `TIME`, `UUID`, `DURATION`, `JSON`. -->

### Read-Only ###

* `id` - (string) The unique identifier (id) of the User.
* `last_login` - (string) The last time the User logged into the application.
* `auth_provider` - (string) The authentication provider for the User.
* `metadata` - (map) A map of the Metadata of the User:
  * `created_at` - (string) The timestamp when this object was created in [RFC3339](https://datatracker.ietf.org/doc/html/rfc3339#section-5.8) format.
  * `created_by` - (string) The user that created this object.
  * `modified_at` - (string) The timestamp when this object was last modified in [RFC3339](https://datatracker.ietf.org/doc/html/rfc3339#section-5.8) format.
  * `modified_by` - (string) The user that modified this object last.
  * `revision_id` - (string) An integer that represent the current revision of the object.

## Importing

An existing User can be [imported](https://www.terraform.io/docs/import/index.html) into this resource using the following command:

```bash
terraform import hyperfabric_user.example_user {userId|email}
```

Starting in Terraform version 1.5, an existing User can be imported
using [import blocks](https://developer.hashicorp.com/terraform/language/import) via the following configuration:

```hcl
import {
  id = "{userId|email}"
  to = hyperfabric_user.example_user
}
```

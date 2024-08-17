---
subcategory: "Administration"
layout: "hyperfabric"
page_title: "Nexus Hyperfabric: hyperfabric_user"
sidebar_current: "docs-hyperfabric-data-source-hyperfabric_user"
description: |-
  Data source for a Nexus Hyperfabric User
---

# hyperfabric_user

Data source for a Nexus Hyperfabric User

A User is a Cisco.com account authorized to access a specific Cisco Nexus Hyperfabric Organization with a specific role that represents the level of privilege of the User.

## API Paths ##

* `/users/{userId|email}` `GET`

## GUI Information ##

* Location: `> Administration > User management`

## Example Usage ##

```hcl
data "hyperfabric_user" "example_user" {
  email = "my-example-user@mydomain.mytld"
}
```

## Schema ##

### Required ###

* `email` - (string) The email of the User.

### Read-Only ###

* `id` - (string) The unique identifier (id) of the User.
* `auth_provider` - (string) The authentication provider for the User.
* `enabled` - (bool) The enabled state of the User.
* `last_login` - (string) The last time the User logged into the application.
* `metadata` - (map) A map of the Metadata of the User:
  * `created_at` - (string) The timestamp when this object was created in [RFC3339](https://datatracker.ietf.org/doc/html/rfc3339#section-5.8) format.
  * `created_by` - (string) The user that created this object.
  * `modified_at` - (string) The timestamp when this object was last modified in [RFC3339](https://datatracker.ietf.org/doc/html/rfc3339#section-5.8) format.
  * `modified_by` - (string) The user that modified this object last.
  * `revision_id` - (string) An integer that represent the current revision of the object.
* `role` - (string) The role assigned to the User that represents the level of privilege of the User.
  - Possible Values: `ADMIN`, `READ_WRITE`, `READ_ONLY`.
* `labels` - (list of strings) A list of user-defined labels that can be used for grouping and filtering objects.
<!-- * `annotations` - (list of maps) A list of key-value annotations to store user-defined data including complex data such as JSON.
  * `name` - (string) The name used to uniquely identify the annotation.
  * `value` - (string) The value of the annotation.
  * `data_type` - (string) The type of data stored in the value of the annotation.
      - Possible Values: `STRING`, `INT32`, `UINT32`, `INT64`, `UINT64`, `BOOL`, `TIME`, `UUID`, `DURATION`, `JSON`. -->

---
subcategory: "Blueprint"
layout: "hyperfabric"
page_title: "Nexus Hyperfabric: hyperfabric_fabric"
sidebar_current: "docs-hyperfabric-resource-hyperfabric_fabric"
description: |-
  Manages a Nexus Hyperfabric Fabric
---

# hyperfabric_fabric

Manages a Nexus Hyperfabric Fabric.

A Fabric is a collection of Nodes, Connections that represents the interconnections between the Nodes, the configuration of the Ports of the Nodes and the logical constructs deployed across the Fabric such as VRFs, logical networks named VNIs and other services.

## API Paths ##

* `/fabrics` `POST`
* `/fabrics/{fabricId|name}` `GET, PUT, DELETE`

## GUI Information ##

* Location: `> Fabrics`

## Example Usage ##

The configuration snippet below creates a Fabric with only the required attributes.

```hcl
resource "hyperfabric_fabric" "example_fabric" {
  name = "my-example-fabric"
}
```
The configuration snippet below shows all possible attributes of a Fabric.

```hcl
resource "hyperfabric_fabric" "full_example_fabric" {
  name        = "my-full-example-fabric"
  description = "This fabric is powered by Cisco Nexus Hyperfabric"
  address     = "170 West Tasman Dr."
  city        = "San Jose"
  country     = "USA"
  location    = "sj01-1-101-AAA01"
  labels = [
    "sj01-1-101-AAA01",
    "blue"
  ]
  annotations = [
    {
      data_type = "STRING"
      name      = "color"
      value     = "blue"
    },
    {
      name  = "rack"
      value = "AAA01"
    }
  ]
}
```

## Schema ##

### Required ###

* `name` - (string) The name of the Fabric.

### Optional ###

* `description` - (string) The description is a user defined field to store notes about the Fabric.
* `topology` - (string) The topology type of the Fabric.
    - Default: `MESH`
    - Valid Values: `MESH`, `SPINE_LEAF`.
* `address` - (string) The physical street address where the Fabric is located.
* `city` - (string) The city in which the Fabric is located.
* `country` - (string) The country in which the Fabric is located.
* `location` - (string) The location is a user defined location of the Fabric.
* `labels` - (list of strings) A list of user-defined labels that can be used for grouping and filtering Fabrics.
* `annotations` - (list of maps) A list of key-value annotations to store user-defined data including complex data such as JSON.


  #### Required ####
  * `name` - (string) The name used to uniquely identify the annotation.
  * `value` - (string) The value of the annotation.


  #### Optional ####
  * `data_type` - (string) The type of data stored in the value of the annotation.
      - Default: `STRING`
      - Valid Values: `STRING`, `INT32`, `UINT32`, `INT64`, `UINT64`, `BOOL`, `TIME`, `UUID`, `DURATION`, `JSON`.

### Read-Only ###

* `id` - (string) The unique identifier (id) of the Fabric.
* `metadata` - (map) A map of the Metadata of the Fabric:
  * `created_at` - (string) The timestamp when this object was created in [RFC3339](https://datatracker.ietf.org/doc/html/rfc3339#section-5.8) format.
  * `created_by` - (string) The user that created this object.
  * `modified_at` - (string) The timestamp when this object was last modified in [RFC3339](https://datatracker.ietf.org/doc/html/rfc3339#section-5.8) format.
  * `modified_by` - (string) The user that modified this object last.
  * `revision_id` - (string) An integer that represent the current revision of the object.

## Importing

An existing Fabric can be [imported](https://www.terraform.io/docs/import/index.html) into this resource using the following command:

```bash
terraform import hyperfabric_fabric.example_fabric {fabricId|name}
```

Starting in Terraform version 1.5, an existing Fabric can be imported
using [import blocks](https://developer.hashicorp.com/terraform/language/import) via the following configuration:

```hcl
import {
  id = "{fabricId|name}"
  to = hyperfabric_fabric.example_fabric
}
```

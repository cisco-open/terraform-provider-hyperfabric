---
subcategory: "Blueprint"
layout: "hyperfabric"
page_title: "Nexus Hyperfabric: hyperfabric_fabric"
sidebar_current: "docs-hyperfabric-data-source-hyperfabric_fabric"
description: |-
  Data source for a Nexus Hyperfabric Fabric
---

# hyperfabric_fabric

Data source for a Nexus Hyperfabric Fabric

A Fabric is a collection of Nodes, Connections that represents the interconnections between the Nodes, the configuration of the Ports of the Nodes and the logical constructs deployed across the Fabric such as VRFs, logical networks named VNIs and other services.

## API Paths ##

* `/fabrics/{fabricId|name}` `GET`

## GUI Information ##

* Location: `> Fabrics`

## Example Usage ##

```hcl
data "hyperfabric_fabric" "example_fabric" {
  name = "my-example-fabric"
}
```

## Schema ##

### Required ###

* `name` - (string) The name of the Fabric.

### Read-Only ###

* `id` - (string) The unique identifier (id) of the Fabric.
* `metadata` - (map) A map of the Metadata of the Fabric:
  * `created_at` - (string) The timestamp when this object was created in [RFC3339](https://datatracker.ietf.org/doc/html/rfc3339#section-5.8) format.
  * `created_by` - (string) The user that created this object.
  * `modified_at` - (string) The timestamp when this object was last modified in [RFC3339](https://datatracker.ietf.org/doc/html/rfc3339#section-5.8) format.
  * `modified_by` - (string) The user that modified this object last.
  * `revision_id` - (string) An integer that represent the current revision of the object.
* `description` - (string) The description is a user defined field to store notes about the Fabric.
* `address` - (string) The physical street address where the Fabric is located.
* `city` - (string) The city in which the Fabric is located.
* `country` - (string) The country in which the Fabric is located.
* `location` - (string) The location is a user defined location of the Fabric.
* `labels` - (list of strings) A list of user-defined labels that can be used for grouping and filtering Fabrics.
* `annotations` - (list of maps) A list of key-value annotations to store user-defined data including complex data such as JSON.
  * `name` - (string) The name used to uniquely identify the annotation.
  * `value` - (string) The value of the annotation.
  * `data_type` - (string) The type of data stored in the value of the annotation.
      - Possible Values: `STRING`, `INT32`, `UINT32`, `INT64`, `UINT64`, `BOOL`, `TIME`, `UUID`, `DURATION`, `JSON`.

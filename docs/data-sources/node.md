---
subcategory: "Blueprint"
layout: "hyperfabric"
page_title: "Nexus Hyperfabric: hyperfabric_node"
sidebar_current: "docs-hyperfabric-data-source-hyperfabric_node"
description: |-
  Data source for a Node in a Nexus Hyperfabric Fabric
---

# hyperfabric_node

Data source for a Node in a Nexus Hyperfabric Fabric

A Node is a logical representation of a device in a Fabric that allows the separation of the logical configuration from the actual physical Device simplifying RMA and hardware replacements. When associated or bound to a Node, a Device assumes the Node identity and all its associated configuration. A Node can be pre-configured and referenced in other Fabric level constructs such as VRFs, VNIs and Link Aggregation Groups (LAGs) before a Device is bound to it allowing for pre-configuration of a complete Fabric.

## API Paths ##

* `/fabrics/{fabricId|fabricName}/nodes/{nodeId|name}` `GET`

## GUI Information ##

* Location: `> Fabrics > {fabric}`

## Example Usage ##

```hcl
data "hyperfabric_node" "example_node" {
  fabric_id = hyperfabric_fabric.example_fabric.id
  name = "my-example-node"
}
```

## Schema ##

### Required ###
* `fabric_id` - (string) The unique identifier (id) of the Fabric. Use the id attribute of the [hyperfabric_fabric](https://registry.terraform.io/providers/cisco-open/hyperfabric/latest/docs/resources/fabric) resource or [hyperfabric_fabric](https://registry.terraform.io/providers/cisco-open/hyperfabric/latest/docs/data-sources/fabric) data source.
* `name` - (string) The name of the Node. The name is used as hostname for the Node and need to comply with DNS restrictions and must be unique in the Fabric.

### Read-Only ###

* `id` - (string) The unique identifier (id) of the Node in the Fabric.
* `position` - (string) The topological position of the Node in the Fabric.
* `device_id` - (string) The unique identifier (id) of the Device associated with the Node.
* `metadata` - (map) A map of the Metadata of the Node:
  * `created_at` - (string) The timestamp when this object was created in [RFC3339](https://datatracker.ietf.org/doc/html/rfc3339#section-5.8) format.
  * `created_by` - (string) The user that created this object.
  * `modified_at` - (string) The timestamp when this object was last modified in [RFC3339](https://datatracker.ietf.org/doc/html/rfc3339#section-5.8) format.
  * `modified_by` - (string) The user that modified this object last.
  * `revision_id` - (string) An integer that represent the current revision of the object.
* `model_name` - (string) The name of model of the Node.
  - Possible Values: `HF6100-32D`, `HF6100-60L4D`.
* `roles` - (list of strings) A list of roles for the Node.
  - Possible Values: `LEAF`, `SPINE`.
* `description` - (string) The description is a user defined field to store notes about the Node.
* `enabled` - (bool) The enabled state of the Node.
* `serial_number` - (string) The serial number of the Device to be associated with the Node.
* `location` - (string) The location is a user defined location of the Node.
* `labels` - (list of strings) A list of user-defined labels that can be used for grouping and filtering objects.
* `annotations` - (list of maps) A list of key-value annotations to store user-defined data including complex data such as JSON.
  * `name` - (string) The name used to uniquely identify the annotation.
  * `value` - (string) The value of the annotation.
  * `data_type` - (string) The type of data stored in the value of the annotation.
      - Possible Values: `STRING`, `INT32`, `UINT32`, `INT64`, `UINT64`, `BOOL`, `TIME`, `UUID`, `DURATION`, `JSON`.

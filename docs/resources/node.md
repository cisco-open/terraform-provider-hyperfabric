---
subcategory: "Blueprint"
layout: "hyperfabric"
page_title: "Nexus Hyperfabric: hyperfabric_node"
sidebar_current: "docs-hyperfabric-resource-hyperfabric_node"
description: |-
  Manages a Node in a Nexus Hyperfabric Fabric
---

# hyperfabric_node

Manages a Node in a Nexus Hyperfabric Fabric

A Node is a logical representation of a device in a Fabric that allows the separation of the logical configuration from the actual physical Device simplifying RMA and hardware replacements. When associated or bound to a Node, a Device assumes the Node identity and all its associated configuration. A Node can be pre-configured and referenced in other Fabric level constructs such as VRFs, VNIs and Link Aggregation Groups (LAGs) before a Device is bound to it allowing for pre-configuration of a complete Fabric.

## API Paths ##

* `/fabrics/{fabricId|fabricName}/nodes` `POST`
* `/fabrics/{fabricId|fabricName}/nodes/{nodeId|name}` `GET, PUT, DELETE`

## GUI Information ##

* Location: `> Fabrics > {fabric}`

## Example Usage ##

The configuration snippet below creates a Node with only the required attributes.

```hcl
resource "hyperfabric_node" "example_node" {
  fabric_id = hyperfabric_fabric.example_fabric.id
  name = "my-example-node"
  model_name = "HF6100-32D"
  roles = ["LEAF"]
}
```
The configuration snippet below shows all possible attributes of a Node.

```hcl
resource "hyperfabric_node" "full_example_node" {
  fabric_id = hyperfabric_fabric.example_fabric.id
  name = "my-full-example-node"
  description = "This node is part of a Cisco Nexus Hyperfabric"
  model_name = "HF6100-32D"
  roles = ["LEAF"]
  enabled = true
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
* `fabric_id` - (string) The unique identifier (id) of the Fabric. Use the id attribute of the [hyperfabric_fabric](https://registry.terraform.io/providers/cisco-open/hyperfabric/latest/docs/resources/fabric) resource or [hyperfabric_fabric](https://registry.terraform.io/providers/cisco-open/hyperfabric/latest/docs/data-sources/fabric) data source.
* `name` - (string) The name of the Node. The name is used as hostname for the Node and need to comply with DNS restrictions and must be unique in the Fabric.
* `model_name` - (string) The name of model of the Node.
  - Valid Values: `HF6100-32D`, `HF6100-60L4D`.
* `roles` - (list of strings) A list of roles for the Node.
  - Valid Values: `LEAF`, `SPINE`.

### Optional ###
  

* `description` - (string) The description is a user defined field to store notes about the Node.
* `enabled` - (bool) The enabled state of the Node.
  - Default: `true`
* `serial_number` - (string) The serial number of Device to be associated with the Node.
* `location` - (string) The location is a user defined location of the Node.
* `labels` - (list of strings) A list of user-defined labels that can be used for grouping and filtering objects.
* `annotations` - (list of maps) A list of key-value annotations to store user-defined data including complex data such as JSON.

  #### Required ####

  * `name` - (string) The name used to uniquely identify the annotation.
  * `value` - (string) The value of the annotation.

  #### Optional ####

  * `data_type` - (string) The type of data stored in the value of the annotation.
      - Default: `STRING`
      - Valid Values: `STRING`, `INT32`, `UINT32`, `INT64`, `UINT64`, `BOOL`, `TIME`, `UUID`, `DURATION`, `JSON`.

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

## Importing

An existing Node can be [imported](https://www.terraform.io/docs/import/index.html) into this resource using the following command:

```bash
terraform import hyperfabric_node.example_node {fabricId|fabricName}/nodes/{nodeId|name}
```

Starting in Terraform version 1.5, an existing Node can be imported
using [import blocks](https://developer.hashicorp.com/terraform/language/import) via the following configuration:

```hcl
import {
  id = "{fabricId|fabricName}/nodes/{nodeId|name}"
  to = hyperfabric_node.example_node
}
```

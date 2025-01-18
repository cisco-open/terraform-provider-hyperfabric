---
subcategory: "Blueprint"
layout: "hyperfabric"
page_title: "Nexus Hyperfabric: hyperfabric_connection"
sidebar_current: "docs-hyperfabric-resource-hyperfabric_connection"
description: |-
  Manages a Connection between two Nodes in a Nexus Hyperfabric Fabric
---

# hyperfabric_connection

Manages a Connection between two Nodes in a Nexus Hyperfabric Fabric

A Connection represents the interconnection between two Ports of two Nodes in a Fabric. Cisco Nexus Hyperfabric uses the connections to generate a possible Bill Of Material, cabling plan and to verify the correct implementation of the desired connectivity intent in a Fabric.

## API Paths ##

* `/fabrics/{fabricId|fabricName}/connections` `POST`
* `/fabrics/{fabricId|fabricName}/connections/{connectionId}` `GET, PUT, DELETE`

## GUI Information ##

* Location: `> Fabrics > {fabric} > Port connections`

## Example Usage ##

The configuration snippet below creates a Connection with only the required attributes.

```hcl
resource "hyperfabric_connection" "example_connection" {
  fabric_id = hyperfabric_fabric.example_fabric.id
  local = {
    node_id = hyperfabric_node.example_node1.node_id
    port_name = "Ethernet1_1"
  }
  remote = {
    node_id = hyperfabric_node.example_node2.node_id
    port_name = "Ethernet1_1"
  }
}
```
The configuration snippet below shows all possible attributes of a Connection.

```hcl
resource "hyperfabric_connection" "full_example_connection" {
  fabric_id = hyperfabric_fabric.example_fabric.id
  description = "This connection is part of a Cisco Nexus Hyperfabric"
  local = {
    node_id = hyperfabric_node.example_node1.node_id
    port_name = "Ethernet1_1"
  }
  remote = {
    node_id = hyperfabric_node.example_node2.node_id
    port_name = "Ethernet1_1"
  }
  pluggable = "QDD-400-AOC7M"
}
```

## Schema ##

### Required ###
* `fabric_id` - (string) The unique identifier (id) of the Fabric. Use the id attribute of the [hyperfabric_fabric](https://registry.terraform.io/providers/cisco-open/hyperfabric/latest/docs/resources/fabric) resource or [hyperfabric_fabric](https://registry.terraform.io/providers/cisco-open/hyperfabric/latest/docs/data-sources/fabric) data source.
* `local` - (map) A map that represents the local side of the Connection.

  #### Required ####

  * `node_id` - (string) The Node unique identifier (node_id) of a Node used as local side of this Connection. Use the node_id attribute of the [hyperfabric_node](https://registry.terraform.io/providers/cisco-open/hyperfabric/latest/docs/resources/node) resource or the [hyperfabric_node](https://registry.terraform.io/providers/cisco-open/hyperfabric/latest/docs/data-sources/node) data source.
  * `port_name` - (string) The name of the Port on the Node used as local side of this Connection.

  #### Read-Only ####

  * `node_name` - (string) The name of the referenced Node used as local side of this Connection.
* `remote` - (map) A map that represents the remote side of the Connection.

  #### Required ####

  * `node_id` - (string) The Node unique identifier (node_id) of a Node used as remote side of this Connection. Use the node_id attribute of the [hyperfabric_node](https://registry.terraform.io/providers/cisco-open/hyperfabric/latest/docs/resources/node) resource or the [hyperfabric_node](https://registry.terraform.io/providers/cisco-open/hyperfabric/latest/docs/data-sources/node) data source.
  * `port_name` - (string) The name of the Port on the Node used as remote side of this Connection.

  #### Read-Only ####

  * `node_name` - (string) The name of the referenced Node used as remote side of this Connection.

### Optional ###
  

* `description` - (string) The description is a user defined field to store notes about the Connection.
<!-- * `cable_type` - (string) The type of cable used for the Connection.
  - Valid Values: `DAC`, `FIBER`.
* `cable_length` - (string) The length of the cable used for the Connection. -->
* `pluggable` - (string) The type of pluggable used for the Connection.
<!-- * `labels` - (list of strings) A list of user-defined labels that can be used for grouping and filtering objects.
* `annotations` - (list of maps) A list of key-value annotations to store user-defined data including complex data such as JSON.

  #### Required ####

  * `name` - (string) The name used to uniquely identify the annotation.
  * `value` - (string) The value of the annotation.

  #### Optional ####

  * `data_type` - (string) The type of data stored in the value of the annotation.
      - Default: `STRING`
      - Valid Values: `STRING`, `INT32`, `UINT32`, `INT64`, `UINT64`, `BOOL`, `TIME`, `UUID`, `DURATION`, `JSON`. -->

### Read-Only ###

* `id` - (string) The unique identifier (id) of the Connection in the Fabric.
* `os_type` - (string) The operating system type of the remote side of the Connection.
* `unrecognized` - (bool) If the remote side of the Connection is recognized or not.
<!-- * `metadata` - (map) A map of the Metadata of the Connection:
  * `created_at` - (string) The timestamp when this object was created in [RFC3339](https://datatracker.ietf.org/doc/html/rfc3339#section-5.8) format.
  * `created_by` - (string) The user that created this object.
  * `modified_at` - (string) The timestamp when this object was last modified in [RFC3339](https://datatracker.ietf.org/doc/html/rfc3339#section-5.8) format.
  * `modified_by` - (string) The user that modified this object last.
  * `revision_id` - (string) An integer that represent the current revision of the object. -->

## Importing

An existing Connection can be [imported](https://www.terraform.io/docs/import/index.html) into this resource using the following command:

```bash
terraform import hyperfabric_connection.example_connection {fabricId|fabricName}/connections/{connectionId}
```

Starting in Terraform version 1.5, an existing Connection can be imported
using [import blocks](https://developer.hashicorp.com/terraform/language/import) via the following configuration:

```hcl
import {
  id = "{fabricId|fabricName}/connections/{connectionId}"
  to = hyperfabric_connection.example_connection
}
```

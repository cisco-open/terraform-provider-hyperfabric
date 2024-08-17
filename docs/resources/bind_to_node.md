---
subcategory: "Blueprint"
layout: "hyperfabric"
page_title: "Nexus Hyperfabric: hyperfabric_bin_to_node"
sidebar_current: "docs-hyperfabric-resource-hyperfabric_bin_to_node"
description: |-
  Manages the binding of a Device to a Node in a Nexus Hyperfabric Fabric
---

# hyperfabric_bin_to_node

Manages the binding of a Device to a Node in a Nexus Hyperfabric Fabric

## API Paths ##

* `/fabrics/{fabricId|fabricName}/nodes/{nodeId|nodeName}` `GET`
* `/fabrics/{fabricId|fabricName}/nodes/{nodeId|nodeName}/devices` `DELETE`
* `/fabrics/{fabricId|fabricName}/nodes/{nodeId|nodeName}/devices/{deviceId}` `PUT`

## GUI Information ##

* Location: `> Fabrics > {fabric} > Nodes > {node} > Bind`

## Example Usage ##

The configuration snippet below binds a Device to a Node.

```hcl
resource "hyperfabric_bin_to_node" "example_bin_to_node" {
  node_id = hyperfabric_node.example_node.id
  device_id = hyperfabric_device.example_device.id
}
```

## Schema ##

### Required ###
* `node_id` - (string) The unique identifier (id) of a Node in a Fabric. Use the id attribute of the [hyperfabric_node](https://registry.terraform.io/providers/cisco-open/hyperfabric/latest/docs/resources/node) resource or [hyperfabric_node](https://registry.terraform.io/providers/cisco-open/hyperfabric/latest/docs/data-sources/node) data source.
* `device_id` - (string) The unique identifier (id) of a Device in a Fabric. Use the id attribute of the [hyperfabric_device](https://registry.terraform.io/providers/cisco-open/hyperfabric/latest/docs/resources/device) resource or [hyperfabric_device](https://registry.terraform.io/providers/cisco-open/hyperfabric/latest/docs/data-sources/device) data source.

### Read-Only ###

* `id` - (string) The unique identifier (id) of the Node in the Fabric.

## Importing

An existing bound Device to a Node can be [imported](https://www.terraform.io/docs/import/index.html) into this resource using the following command:

```bash
terraform import hyperfabric_bin_to_node.example_bin_to_node {fabricId|fabricName}/nodes/{nodeId|nodeName}
```

Starting in Terraform version 1.5, an existing bound Device to a Node can be imported
using [import blocks](https://developer.hashicorp.com/terraform/language/import) via the following configuration:

```hcl
import {
  id = "{fabricId|fabricName}/nodes/{nodeId|nodeName}"
  to = hyperfabric_bin_to_node.example_bin_to_node
}
```

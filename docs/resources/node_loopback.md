---
subcategory: "Blueprint"
layout: "hyperfabric"
page_title: "Nexus Hyperfabric: hyperfabric_node_loopback"
sidebar_current: "docs-hyperfabric-resource-hyperfabric_node_loopback"
description: |-
  Manages a Loopback of a Node in a Nexus Hyperfabric Fabric
---

# hyperfabric_node_loopback

Manages a Loopback of a Node in a Nexus Hyperfabric Fabric

A Loopback is a logical network interface of a Node used as source interface to peer with external devices.

## API Paths ##

* `/fabrics/{fabricId|fabricName}/nodes/{nodeId|nodeName}/loopbacks` `POST`
* `/fabrics/{fabricId|fabricName}/nodes/{nodeId|nodeName}/loopbacks/{loopbackId|name}` `GET, PUT, DELETE`

## GUI Information ##

* Location: `> Fabrics > {fabric} > Route tables (VRF) > {vrf} > Loopback interfaces`

## Example Usage ##

The configuration snippet below creates a Loopback of a Node with only the required attributes.

```hcl
resource "hyperfabric_node_loopback" "example_node_loopback" {
  node_id = hyperfabric_node.example_node.id
  name = "Loopback10"
  ipv4_address = "10.1.0.1/24"
}
```
Or
```hcl
resource "hyperfabric_node_loopback" "example_node_loopback" {
  node_id = hyperfabric_node.example_node.id
  name = "Loopback10"
  ipv6_address = "2001:1::1"
}
```

The configuration snippet below shows all possible attributes of a Loopback of a Node.

```hcl
resource "hyperfabric_node_loopback" "full_example_node_loopback" {
  node_id      = hyperfabric_node.example_node.id
  name         = "Loopback10"
  description  = "Used for BGP peering"
  ipv4_address = "10.1.0.1"
  ipv6_address = "2001:1::1"
  vrf_id       = hyperfabric_vrf.example_vrf.vrf_id
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
      name  = "community"
      value = "AAA01"
    }
  ]
}
```

## Schema ##

### Required ###
* `node_id` - (string) The unique identifier (id) of a Node in a Fabric. Use the id attribute of the [hyperfabric_node](https://registry.terraform.io/providers/cisco-open/hyperfabric/latest/docs/resources/node) resource or [hyperfabric_node](https://registry.terraform.io/providers/cisco-open/hyperfabric/latest/docs/data-sources/node) data source.
* `name` - (string) The name of the Loopback of the Node.

### Optional ###

* `description` - (string) The description is a user defined field to store notes about the Loopback of the Node.
* `ipv4_address` - (string) An IPv4 address without a subnet mask to be configured on the Loopback. One of `ipv4_address` or `ipv6_address` is required.
* `ipv6_address` - (string) An IPv6 address without a subnet mask to be configured on the Loopback. One of `ipv4_address` or `ipv6_address` is required.
* `vrf_id` - (string) The `vrf_id` to associate with the Loopback of the Node. Use the vrf_id attribute of the [hyperfabric_vrf](https://registry.terraform.io/providers/cisco-open/hyperfabric/latest/docs/resources/vrf) resource or [hyperfabric_vrf](https://registry.terraform.io/providers/cisco-open/hyperfabric/latest/docs/data-sources/vrf) data source.
  - Default to the id of the Default-VRF.
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

* `id` - (string) The unique identifier (id) of the Loopback of the Node in the Fabric.
* `metadata` - (map) A map of the Metadata of the Node:
  * `created_at` - (string) The timestamp when this object was created in [RFC3339](https://datatracker.ietf.org/doc/html/rfc3339#section-5.8) format.
  * `created_by` - (string) The user that created this object.
  * `modified_at` - (string) The timestamp when this object was last modified in [RFC3339](https://datatracker.ietf.org/doc/html/rfc3339#section-5.8) format.
  * `modified_by` - (string) The user that modified this object last.
  * `revision_id` - (string) An integer that represent the current revision of the object.

## Importing

An existing Loopback of a Node can be [imported](https://www.terraform.io/docs/import/index.html) into this resource using the following command:

```bash
terraform import hyperfabric_node_loopback.example_node_loopback {fabricId|fabricName}/nodes/{nodeId|nodeName}/loopbacks/{loopbackId|name}
```

Starting in Terraform version 1.5, an existing Loopback of a Node can be imported
using [import blocks](https://developer.hashicorp.com/terraform/language/import) via the following configuration:

```hcl
import {
  id = "{fabricId|fabricName}/nodes/{nodeId|nodeName}/loopbacks/{loopbackId|name}"
  to = hyperfabric_node_loopback.example_node_loopback
}
```

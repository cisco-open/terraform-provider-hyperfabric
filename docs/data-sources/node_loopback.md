---
subcategory: "Blueprint"
layout: "hyperfabric"
page_title: "Nexus Hyperfabric: hyperfabric_node_loopback"
sidebar_current: "docs-hyperfabric-data-source-hyperfabric_node_loopback"
description: |-
  Data source for a Loopback of a Node in a Nexus Hyperfabric Fabric
---

# hyperfabric_node_loopback

Data source for a Loopback of a Node in a Nexus Hyperfabric Fabric

A Loopback is a logical network interface of a Node used as source interface to peer with external devices.

## API Paths ##

* `/fabrics/{fabricId|fabricName}/nodes/{nodeId|nodeName}/loopbacks/{loopbackId|name}` `GET`

## GUI Information ##

* Location: `> Fabrics > {fabric} > Route tables (VRF) > {vrf} > Loopback interfaces`

## Example Usage ##

```hcl
data "hyperfabric_node_loopback" "example_node_loopback" {
  node_id = hyperfabric_node.example_node.id
  name = "Loopback10"
}
```

## Schema ##

### Required ###
* `node_id` - (string) The unique identifier (id) of a Node in a Fabric. Use the id attribute of the [hyperfabric_node](https://registry.terraform.io/providers/cisco-open/hyperfabric/latest/docs/resources/node) resource or [hyperfabric_node](https://registry.terraform.io/providers/cisco-open/hyperfabric/latest/docs/data-sources/node) data source.
* `name` - (string) The name of the Loopback of the Node.

### Read-Only ###

* `id` - (string) The unique identifier (id) of the Loopback of the Node in the Fabric.
* `description` - (string) The description is a user defined field to store notes about the Loopback of the Node.
* `ipv4_address` - (string) The IPv4 address configured on the Loopback.
* `ipv6_address` - (string) The IPv6 address configured on the Loopback.
* `vrf_id` - (string) The `vrf_id` to associate with the Loopback of the Node.
* `metadata` - (map) A map of the Metadata of the Node:
  * `created_at` - (string) The timestamp when this object was created in [RFC3339](https://datatracker.ietf.org/doc/html/rfc3339#section-5.8) format.
  * `created_by` - (string) The user that created this object.
  * `modified_at` - (string) The timestamp when this object was last modified in [RFC3339](https://datatracker.ietf.org/doc/html/rfc3339#section-5.8) format.
  * `modified_by` - (string) The user that modified this object last.
  * `revision_id` - (string) An integer that represent the current revision of the object.
* `labels` - (list of strings) A list of user-defined labels that can be used for grouping and filtering objects.
* `annotations` - (list of maps) A list of key-value annotations to store user-defined data including complex data such as JSON.
  * `name` - (string) The name used to uniquely identify the annotation.
  * `value` - (string) The value of the annotation.
  * `data_type` - (string) The type of data stored in the value of the annotation.
      - Possible Values: `STRING`, `INT32`, `UINT32`, `INT64`, `UINT64`, `BOOL`, `TIME`, `UUID`, `DURATION`, `JSON`.

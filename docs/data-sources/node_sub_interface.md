---
subcategory: "Blueprint"
layout: "hyperfabric"
page_title: "Nexus Hyperfabric: hyperfabric_node_sub_interface"
sidebar_current: "docs-hyperfabric-data-source-hyperfabric_node_sub_interface"
description: |-
  Data source for a Sub-Interface of a Node in a Nexus Hyperfabric Fabric
---

# hyperfabric_node_sub_interface

Data source for a Sub-Interface of a Node in a Nexus Hyperfabric Fabric

A dot1q VLAN Sub-Interface is a virtual network interface that is associated with a VLAN ID on a routed front panel network interface of a Node used as source interface to peer with external devices. The parent interface of a Sub-Interface is a Node Port with the unused or routed role.

## API Paths ##

* `/fabrics/{fabricId|fabricName}/nodes/{nodeId|nodeName}/subInterfaces/{subInterfaceId|name}` `GET`

## GUI Information ##

* Location: `> Fabrics > {fabric} > Nodes > {node} > Configure > Port configuration`

## Example Usage ##

```hcl
data "hyperfabric_node_sub_interface" "example_node_sub_interface" {
  node_id = hyperfabric_node.example_node.id
  name = "Ethernet1_1.100"
}
```

## Schema ##

### Required ###
* `node_id` - (string) The unique identifier (id) of a Node in a Fabric. Use the id attribute of the [hyperfabric_node](https://registry.terraform.io/providers/cisco-open/hyperfabric/latest/docs/resources/node) resource or [hyperfabric_node](https://registry.terraform.io/providers/cisco-open/hyperfabric/latest/docs/data-sources/node) data source.
* `name` - (string) The name of the Sub-Interface of the Node.

### Read-Only ###

* `id` - (string) The unique identifier (id) of the Sub-Interface of the Node in the Fabric.
* `sub_interface_id` - (string) The unique identifier (id) of a Sub-Interface.
* `description` - (string) The description is a user defined field to store notes about the Sub-Interface of the Node.
* `enabled` - (bool) The enabled state of the Sub-Interface of the Node.
* `ipv4_addresses` - (list of strings) A list of IPv4 addresses with subnet mask configured on the Sub-Interface.
* `ipv6_addresses` - (list of strings) A list of IPv6 addresses with subnet mask configured on the Sub-Interface.
* `parent` - (string) The `parent` Port of the Sub-Interface of the Node.
* `vlan_id` - (integer) The VLAN ID used for the encapsulation of the Sub-Interface of the Node.
* `vrf_id` - (string) The `vrf_id` of a VRF to associate with the Sub-Interface of the Node.
* `metadata` - (map) A map of the Metadata of the Node Sub-Interface:
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

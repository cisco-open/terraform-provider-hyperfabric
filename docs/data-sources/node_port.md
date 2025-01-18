---
subcategory: "Blueprint"
layout: "hyperfabric"
page_title: "Nexus Hyperfabric: hyperfabric_node_port"
sidebar_current: "docs-hyperfabric-data-source-hyperfabric_node_port"
description: |-
  Data source for a Port of a Node in a Nexus Hyperfabric Fabric
---

# hyperfabric_node_port

Data source for a Port of a Node in a Nexus Hyperfabric Fabric

A Port is a front panel network interface of a Node used as Fabric Port to interconnect with other Nodes, as Routed Port to peer at Layer 3 with external devices or as a Host Port to connect to other endpoints via Layer 2 (VLAN).

## API Paths ##

* `/fabrics/{fabricId|fabricName}/nodes/{nodeId|nodeName}/ports/{portId|name}` `GET`

## GUI Information ##

* Location: `> Fabrics > {fabric} > Nodes > {node} > Configure > Port configuration`

## Example Usage ##

```hcl
data "hyperfabric_node_port" "example_node_port" {
  node_id = hyperfabric_node.example_node.id
  name = "Ethernet1_1"
}
```

## Schema ##

### Required ###
* `node_id` - (string) The unique identifier (id) of a Node in a Fabric. Use the id attribute of the [hyperfabric_node](https://registry.terraform.io/providers/cisco-open/hyperfabric/latest/docs/resources/node) resource or [hyperfabric_node](https://registry.terraform.io/providers/cisco-open/hyperfabric/latest/docs/data-sources/node) data source.
* `name` - (string) The name of the Port of the Node.

### Read-Only ###

* `id` - (string) The unique identifier (id) of the Port of the Node in the Fabric.
* `index` - (integer) The index number of the Port of the Node.
* `linecard` - (integer) The linecard index number of the Port of the Node.
* `description` - (string) The description is a user defined field to store notes about the Port of the Node.
* `enabled` - (bool) The enabled state of the Port of the Node.
* `ipv4_addresses` - (list of strings) A list of IPv4 addresses with subnet mask configured on the Port.
* `ipv6_addresses` - (list of strings) A list of IPv6 addresses with subnet mask configured on the Port.
* `lldp_host` - (string) The name of host reported by LLDP connected to the Port of the Node.
* `lldp_info` - (string) The info about the host reported by LLDP connected to the Port of the Node.
* `lldp_port` - (string) The port of host reported by LLDP connected to the Port of the Node.
* `max_speed` - (string) The maximum speed of the Port of the Node.
* `mtu` - (integer) The MTU of the Port of the Node.
* `prevent_forwarding` - (bool) Prevent traffic from being forwarded by the Port.
* `roles` - (list of strings) A list of roles configured on the Port.
  - Possible Values: `UNUSED_PORT`, `FABRIC_PORT`, `HOST_PORT`, `ROUTED_PORT`.
* `speed` - (string) The configured speed of the Port of the Node.
* `sub_interfaces_count` - (string) The number of sub-interfaces of the Port of the Node.
* `vlan_ids` - (list of strings) A list of Vlan IDs used by the Port of the Node.
* `vnis` - (list of strings) A list of VNIs used by the Port of the Node.
* `vrf_id` - (string) The `vrf_id` to associate with the Port of the Node.
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

---
subcategory: "Blueprint"
layout: "hyperfabric"
page_title: "Nexus Hyperfabric: hyperfabric_node_port"
sidebar_current: "docs-hyperfabric-resource-hyperfabric_node_port"
description: |-
  Manages a Port of a Node in a Nexus Hyperfabric Fabric
---

# hyperfabric_node_port

Manages a Port of a Node in a Nexus Hyperfabric Fabric

A Port is a front panel network interface of a Node used as Fabric Port to interconnect with other Nodes, as Routed Port to peer at Layer 3 with external devices or as a Host Port to connect to other endpoints via Layer 2 (VLAN).

## API Paths ##

* `/fabrics/{fabricId|fabricName}/nodes/{nodeId|nodeName}/ports/{portId|name}` `GET, PUT, DELETE`

## GUI Information ##

* Location: `> Fabrics > {fabric} > Nodes > {node} > Configure > Port configuration`

## Example Usage ##

The configuration snippet below creates a Port of a Node with only the required attributes.

```hcl
resource "hyperfabric_node_port" "example_node_port" {
  node_id = hyperfabric_node.example_node.id
  name = "Ethernet1_1"
}
```
The configuration snippet below shows all possible attributes of a Port of a Node.

```hcl
resource "hyperfabric_node_port" "full_example_node_port" {
  node_id            = hyperfabric_node.example_node.id
  name               = "Ethernet1_1"
  description        = "Connected to server01"
  enabled            = true
  ipv4_addresses     = ["10.1.0.1/24"]
  ipv6_addresses     = ["2001:1::1/64", "2002:1::1/64"]
  prevent_forwarding = true
  roles              = ["ROUTED_PORT"]
  vrf_id             = hyperfabric_vrf.example_vrf.vrf_id
}
```

## Schema ##

### Required ###
* `node_id` - (string) The unique identifier (id) of a Node in a Fabric. Use the id attribute of the [hyperfabric_node](https://registry.terraform.io/providers/cisco-open/hyperfabric/latest/docs/resources/node) resource or [hyperfabric_node](https://registry.terraform.io/providers/cisco-open/hyperfabric/latest/docs/data-sources/node) data source.
* `name` - (string) The name of the Port of the Node.
* `roles` - (list of strings) A list of roles to be configured on the Port.
  - Valid Values: `UNUSED_PORT`, `FABRIC_PORT`, `HOST_PORT`, `ROUTED_PORT`.

### Optional ###

* `description` - (string) The description is a user defined field to store notes about the Port of the Node.
* `enabled` - (bool) The enabled state of the Port of the Node.
* `ipv4_addresses` - (list of strings) A list of IPv4 addresses with subnet mask to be configured on the Port. Requires the `ROUTED_PORT` role to be configured in `roles` and the `vrf_id` to be set.
* `ipv6_addresses` - (list of strings) A list of IPv6 addresses with subnet mask to be configured on the Port. Requires the `ROUTED_PORT` role to be configured in `roles` and the `vrf_id` to be set.
* `prevent_forwarding` - (bool) Prevent traffic from being forwarded by the Port. Requires `enabled` to be set to `true` (equivalent to `Admin State` set to `Up`) and role to be one of `UNUSED_PORT`, `ROUTED_PORT` or `HOST_PORT`.
* `vrf_id` - (string) The `vrf_id` to associate with the Port of the Node. Use the vrf_id attribute of the [hyperfabric_vrf](https://registry.terraform.io/providers/cisco-open/hyperfabric/latest/docs/resources/vrf) resource or [hyperfabric_vrf](https://registry.terraform.io/providers/cisco-open/hyperfabric/latest/docs/data-sources/vrf) data source.
  - Required when the Port `roles` include `ROUTED_PORT`.
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

* `id` - (string) The unique identifier (id) of a Port of the Node in the Fabric.
* `index` - (integer) The index number of the Port of the Node.
* `linecard` - (integer) The linecard index number of the Port of the Node.
* `lldp_host` - (string) The name of host reported by LLDP connected to the Port of the Node.
* `lldp_info` - (string) The info about the host reported by LLDP connected to the Port of the Node.
* `lldp_port` - (string) The port of host reported by LLDP connected to the Port of the Node.
* `max_speed` - (string) The maximum speed of the Port of the Node.
* `mtu` - (integer) The MTU of the Port of the Node.
* `speed` - (string) The configured speed of the Port of the Node.
* `sub_interfaces_count` - (string) The number of sub-interfaces of the Port of the Node.
* `vlan_ids` - (list of strings) A list of Vlan IDs used by the Port of the Node.
* `vnis` - (list of strings) A list of VNIs used by the Port of the Node.
* `metadata` - (map) A map of the Metadata of the Node:
  * `created_at` - (string) The timestamp when this object was created in [RFC3339](https://datatracker.ietf.org/doc/html/rfc3339#section-5.8) format.
  * `created_by` - (string) The user that created this object.
  * `modified_at` - (string) The timestamp when this object was last modified in [RFC3339](https://datatracker.ietf.org/doc/html/rfc3339#section-5.8) format.
  * `modified_by` - (string) The user that modified this object last.
  * `revision_id` - (string) An integer that represent the current revision of the object.

## Importing

An existing Port of a Node can be [imported](https://www.terraform.io/docs/import/index.html) into this resource using the following command:

```bash
terraform import hyperfabric_node_port.example_node_port {fabricId|fabricName}/nodes/{nodeId|nodeName}/ports/{id|name}
```

Starting in Terraform version 1.5, an existing Port of a Node can be imported
using [import blocks](https://developer.hashicorp.com/terraform/language/import) via the following configuration:

```hcl
import {
  id = "{fabricId|fabricName}/nodes/{nodeId|nodeName}/ports/{id|name}"
  to = hyperfabric_node_port.example_node_port
}
```

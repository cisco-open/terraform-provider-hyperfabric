---
subcategory: "Blueprint"
layout: "hyperfabric"
page_title: "Nexus Hyperfabric: hyperfabric_node_sub_interface"
sidebar_current: "docs-hyperfabric-resource-hyperfabric_node_sub_interface"
description: |-
  Manages a Sub-Interface of a Node in a Nexus Hyperfabric Fabric
---

# hyperfabric_node_sub_interface

Manages a Sub-Interface of a Node in a Nexus Hyperfabric Fabric

A dot1q VLAN Sub-Interface is a virtual network interface that is associated with a VLAN ID on a routed front panel network interface of a Node used as source interface to peer with external devices. The parent interface of a Sub-Interface is a Node Port with the `UNUSED_PORT` or `ROUTED_PORT` role.

## API Paths ##

* `/fabrics/{fabricId|fabricName}/nodes/{nodeId|nodeName}/subInterfaces` `POST`
* `/fabrics/{fabricId|fabricName}/nodes/{nodeId|nodeName}/subInterfaces/{subInterfaceId|name}` `GET, PUT, DELETE`

## GUI Information ##

* Location: `> Fabrics > {fabric} > Nodes > {node} > Configure > Port configuration`

## Example Usage ##

The configuration snippet below creates a Sub-Interface of a Node with only the required attributes.

```hcl
resource "hyperfabric_node_sub_interface" "example_node_sub_interface" {
  node_id = hyperfabric_node.example_node.id
  name = "Ethernet1_1.100"
}
```
The configuration snippet below shows all possible attributes of a Sub-Interface of a Node.

```hcl
resource "hyperfabric_node_sub_interface" "full_example_node_sub_interface" {
  node_id        = hyperfabric_node.example_node.id
  name           = "Ethernet1_1.100"
  description    = "Loopback for BGP peering"
  enabled        = true
  ipv4_addresses = ["10.1.0.1/24"]
  ipv6_addresses = ["2001:1::1/64", "2002:1::1/64"]
	vlan_id        = 100
  vrf_id         = hyperfabric_vrf.example_vrf.vrf_id
	labels         = [
		"sj01-1-101-AAA01",
		"blue"
	]
	annotations    = [
		{
			name      = "color"
			value     = "blue"
		},
		{
			data_type = "UINT32"
			name      = "rack"
			value     = "1"
		}
	]
}
```

## Schema ##

### Required ###
* `node_id` - (string) The unique identifier (id) of a Node in a Fabric. Use the id attribute of the [hyperfabric_node](https://registry.terraform.io/providers/cisco-open/hyperfabric/latest/docs/resources/node) resource or [hyperfabric_node](https://registry.terraform.io/providers/cisco-open/hyperfabric/latest/docs/data-sources/node) data source.
* `name` - (string) The name of the Sub-Interface of the Node. The name should be in the `<Port Name>.<Integer>` format (i.e. `Ethernet1_1.100`). If `vlan_id` attribute is not provided, the integer in the Sub-Interface name will be used as the encapsulation VLAN ID.

### Optional ###

* `description` - (string) The description is a user defined field to store notes about the Sub-Interface of the Node.
* `enabled` - (bool) The enabled state of the Sub-Interface of the Node.
* `ipv4_addresses` - (list of strings) A list of IPv4 addresses with subnet mask to be configured on the Sub-Interface.
* `ipv6_addresses` - (list of strings) A list of IPv6 addresses with subnet mask to be configured on the Sub-Interface.
* `vlan_id` - (integer) The VLAN ID to use as encapsulation for the Sub-Interface of the Node. If not provided, the integer in the Sub-Interface `name` will be used as the encapsulation VLAN ID.
* `vrf_id` - (string) The `vrf_id` of a VRF to associate with the Sub-Interface of the Node. Use the vrf_id attribute of the [hyperfabric_vrf](https://registry.terraform.io/providers/cisco-open/hyperfabric/latest/docs/resources/vrf) resource or [hyperfabric_vrf](https://registry.terraform.io/providers/cisco-open/hyperfabric/latest/docs/data-sources/vrf) data source.
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

* `id` - (string) The unique identifier (id) of a Sub-Interface of the Node in the Fabric.
* `sub_interface_id` - (string) The unique identifier (id) of a Sub-Interface.
* `parent` - (string) The `parent` Port of the Sub-Interface of the Node.
* `metadata` - (map) A map of the Metadata of the Node Sub-Interface:
  * `created_at` - (string) The timestamp when this object was created in [RFC3339](https://datatracker.ietf.org/doc/html/rfc3339#section-5.8) format.
  * `created_by` - (string) The user that created this object.
  * `modified_at` - (string) The timestamp when this object was last modified in [RFC3339](https://datatracker.ietf.org/doc/html/rfc3339#section-5.8) format.
  * `modified_by` - (string) The user that modified this object last.
  * `revision_id` - (string) An integer that represent the current revision of the object.

## Importing

An existing Sub-Interface of a Node can be [imported](https://www.terraform.io/docs/import/index.html) into this resource using the following command:

```bash
terraform import hyperfabric_node_sub_interface.example_node_sub_interface {fabricId|fabricName}/nodes/{nodeId|nodeName}/subInterfaces/{subInterfaceId|name}
```

Starting in Terraform version 1.5, an existing Sub-Interface of a Node can be imported
using [import blocks](https://developer.hashicorp.com/terraform/language/import) via the following configuration:

```hcl
import {
  id = "{fabricId|fabricName}/nodes/{nodeId|nodeName}/subInterfaces/{subInterfaceId|name}"
  to = hyperfabric_node_sub_interface.example_node_sub_interface
}
```

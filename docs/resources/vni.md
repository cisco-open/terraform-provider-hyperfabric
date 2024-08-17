---
subcategory: "Networking"
layout: "hyperfabric"
page_title: "Nexus Hyperfabric: hyperfabric_vni"
sidebar_current: "docs-hyperfabric-resource-hyperfabric_vni"
description: |-
  Manages a VNI in a Nexus Hyperfabric Fabric
---

# hyperfabric_vni

Manages a VNI in a Nexus Hyperfabric Fabric

A VNI represents a Layer 2 or Layer 3 logical network that can be extended across the Fabric and mapped to a VLAN ID on specific Ports and LAGs. A VNI can be mapped to a VRF and configured with an SVI to serve as an Anycast Gateway.

## API Paths ##

* `/fabrics/{fabricId|fabricName}/vnis` `POST`
* `/fabrics/{fabricId|fabricName}/vnis/{vniId|name}` `GET, PUT, DELETE`

## GUI Information ##

* Location:
  - `> Fabrics > {fabric} > Logical network > Logical Networks (VNI)`
  - `> Fabrics > {fabric} > Attachments > VLAN memberships`

## Example Usage ##

The configuration snippet below creates a VNI with only the required attributes.

```hcl
resource "hyperfabric_vni" "example_vni" {
  fabric_id = hyperfabric_fabric.example_fabric.id
  name      = "my-example-vni"
}
```
The configuration snippet below shows all possible attributes of a VNI.

```hcl
resource "hyperfabric_vni" "full_example_vni" {
  fabric_id   = hyperfabric_fabric.example_fabric.id
  name        = "my-full-example-vni"
  description = "This VNI is part of a Cisco Nexus Hyperfabric"
  vni         = 10000
  svi = {
    enabled        = true
    ipv4_addresses = ["192.168.0.254/24"]
    ipv6_addresses = ["2001::1/64", "2002::1/64"]
  }
  members = [
    {
      node_id = "*"
      port_name = "*"
      vlan_id = 103
    },
    {
      node_id = hyperfabric_node.example_node.node_id
      port_name = "*"
      vlan_id   = 103
    },
    {
      node_id   = hyperfabric_node.example_node.node_id
      port_name = "Ethernet1_10"
      vlan_id   = 103
    }
  ]
  vrf_id = hyperfabric_vrf.example_vrf.vrf_id
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
* `name` - (string) The name of the VNI.

### Optional ###

* `description` - (string) The description is a user defined field to store notes about the VNI.
* `vni` - (integer) The VXLAN Network Identifier (VNID) used for the VNI.
* `members` - (list of maps) A list of key-value annotations to store user-defined data including complex data such as JSON.

  #### Required ####

  * `node_id` - (string) The unique identifier (nodeId) of the Node. Use the node_id attribute of the [hyperfabric_node](https://registry.terraform.io/providers/cisco-open/hyperfabric/latest/docs/resources/node) resource or [hyperfabric_node](https://registry.terraform.io/providers/cisco-open/hyperfabric/latest/docs/data-sources/node) data source or "*" for all Nodes.
  * `port_name` - (string) The name of the Port or "*" for all ports on a Node or all Nodes.

  #### Read-Only ####

  * `node_name` - (string) The name of the Node referenced by `node_id` for this member.
* `vrf_id` - (string) The unique identifier (vrfId) of the VRF. Use the vrf_id attribute of the [hyperfabric_vrf](https://registry.terraform.io/providers/cisco-open/hyperfabric/latest/docs/resources/vrf) resource or [hyperfabric_vrf](https://registry.terraform.io/providers/cisco-open/hyperfabric/latest/docs/data-sources/vrf) data source.
* `svi` - (map) A map of the attributes of the SVI for the VNI. Requires `vrf_id` to also be set.

  #### Required At Least One Of ####

  * `ipv4_addresses` - (list of strings) A list of IPv4 addresses with their subnet Mask to be used by the SVI Anycast Gateway.
  * `ipv6_addresses` - (list of strings) A list of IPv6 addresses with their subnet Mask to be used by the SVI Anycast Gateway.

  #### Optional ####

  * `enabled` - (string) The enabled state of the SVI.
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

* `id` - (string) The unique identifier (id) of the VNI in the Fabric.
* `vni_id` - (string) The unique identifier (id) of the VNI.
* `enabled` - (bool) The enabled state of the VNI.
* `is_default` - (string) The flag that denote if the VNI is the default VNI or not.
* `mtu` - (integer) The MTU of the SVI of the VNI.
* `metadata` - (map) A map of the Metadata of the VNI:
  * `created_at` - (string) The timestamp when this object was created in [RFC3339](https://datatracker.ietf.org/doc/html/rfc3339#section-5.8) format.
  * `created_by` - (string) The user that created this object.
  * `modified_at` - (string) The timestamp when this object was last modified in [RFC3339](https://datatracker.ietf.org/doc/html/rfc3339#section-5.8) format.
  * `modified_by` - (string) The user that modified this object last.
  * `revision_id` - (string) An integer that represent the current revision of the object.

## Importing

An existing VNI can be [imported](https://www.terraform.io/docs/import/index.html) into this resource using the following command:

```bash
terraform import hyperfabric_vni.example_vni {fabricId|fabricName}/vnis/{vniId|name}
```

Starting in Terraform version 1.5, an existing VNI can be imported
using [import blocks](https://developer.hashicorp.com/terraform/language/import) via the following configuration:

```hcl
import {
  id = "{fabricId|fabricName}/vnis/{vniId|name}"
  to = hyperfabric_vni.example_vni
}
```

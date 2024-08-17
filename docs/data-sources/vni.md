---
subcategory: "Networking"
layout: "hyperfabric"
page_title: "Nexus Hyperfabric: hyperfabric_vni"
sidebar_current: "docs-hyperfabric-data-source-hyperfabric_vni"
description: |-
  Manages a VNI in a Nexus Hyperfabric Fabric
---

# hyperfabric_vni

Manages a VNI in a Nexus Hyperfabric Fabric

A VNI represents a Layer 2 or Layer 3 logical network that can be extended across the Fabric and mapped to a VLAN ID on specific Ports and LAGs. A VNI can be mapped to a VRF and configured with an SVI to serve as an Anycast Gateway.

## API Paths ##

* `/fabrics/{fabricId|fabricName}/vnis/{vniId|name}` `GET`

## GUI Information ##

* Location:
  - `> Fabrics > {fabric} > Logical network > Logical Networks (VNI)`
  - `> Fabrics > {fabric} > Attachments > VLAN memberships`

## Example Usage ##

```hcl
data "hyperfabric_vni" "example_vni" {
  fabric_id = hyperfabric_fabric.example_fabric.id
  name      = "my-example-vni"
}
```

## Schema ##

### Required ###
* `fabric_id` - (string) The unique identifier (id) of the Fabric. Use the id attribute of the [hyperfabric_fabric](https://registry.terraform.io/providers/cisco-open/hyperfabric/latest/docs/resources/fabric) resource or [hyperfabric_fabric](https://registry.terraform.io/providers/cisco-open/hyperfabric/latest/docs/data-sources/fabric) data source.
* `name` - (string) The name of the VNI.

### Read-Only ###

* `id` - (string) The unique identifier (id) of the VNI in the Fabric.
* `vni_id` - (string) The unique identifier (id) of the VNI.
* `description` - (string) The description is a user defined field to store notes about the VNI.
* `enabled` - (bool) The enabled state of the VNI.
* `is_default` - (string) The flag that denote if the VNI is the default VNI or not.
* `mtu` - (integer) The MTU of the SVI of the VNI.
* `members` - (list of maps) A list of key-value annotations to store user-defined data including complex data such as JSON.
  * `node_id` - (string) The unique identifier (nodeId) of the Node or "*" for all Nodes.
  * `port_name` - (string) The name of the Port or "*" for all ports on a Node or all Nodes.
  * `node_name` - (string) The name of the Node referenced by `node_id` for this member.
* `metadata` - (map) A map of the Metadata of the VNI:
  * `created_at` - (string) The timestamp when this object was created in [RFC3339](https://datatracker.ietf.org/doc/html/rfc3339#section-5.8) format.
  * `created_by` - (string) The user that created this object.
  * `modified_at` - (string) The timestamp when this object was last modified in [RFC3339](https://datatracker.ietf.org/doc/html/rfc3339#section-5.8) format.
  * `modified_by` - (string) The user that modified this object last.
  * `revision_id` - (string) An integer that represent the current revision of the object.
* `svi` - (map) A map of the attributes of the SVI for the VNI.
  * `ipv4_addresses` - (list of strings) A list of IPv4 addresses with their subnet Mask to be used by the SVI Anycast Gateway.
  * `ipv6_addresses` - (list of strings) A list of IPv6 addresses with their subnet Mask to be used by the SVI Anycast Gateway.
  * `enabled` - (string) The enabled state of the SVI.
* `vni` - (integer) The VXLAN Network Identifier (VNID) used for the VNI.
* `vrf_id` - (string) The unique identifier (vrfId) of the VRF.
* `labels` - (list of strings) A list of user-defined labels that can be used for grouping and filtering objects.
* `annotations` - (list of maps) A list of key-value annotations to store user-defined data including complex data such as JSON.
  * `name` - (string) The name used to uniquely identify the annotation.
  * `value` - (string) The value of the annotation.
  * `data_type` - (string) The type of data stored in the value of the annotation.
      - Default: `STRING`
      - Valid Values: `STRING`, `INT32`, `UINT32`, `INT64`, `UINT64`, `BOOL`, `TIME`, `UUID`, `DURATION`, `JSON`.

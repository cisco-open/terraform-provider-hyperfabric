---
subcategory: "Networking"
layout: "hyperfabric"
page_title: "Nexus Hyperfabric: hyperfabric_vrf"
sidebar_current: "docs-hyperfabric-data-source-hyperfabric_vrf"
description: |-
  Data source for a VRF in a Nexus Hyperfabric Fabric
---

# hyperfabric_vrf

Data source for a VRF in a Nexus Hyperfabric Fabric

A VRF is a virtual-routing-and-forwarding instance that represents a routing table deployed across Nodes in the Fabric to implement VRF-lite and the ability to have multiple routing tables on a single Device. A VRF can be associated to multiple VNIs configured with an Anycast Gateway SVI.

## API Paths ##

* `/fabrics/{fabricId|fabricName}/vrfs/{vrfId|name}` `GET`

## GUI Information ##

* Location: `> Fabrics > {fabric} > Logical network > Route tables (VRF)`

## Example Usage ##

```hcl
data "hyperfabric_vrf" "example_vrf" {
  fabric_id = hyperfabric_fabric.example_fabric.id
  name      = "my-example-vrf"
}
```

## Schema ##

### Required ###
* `fabric_id` - (string) The unique identifier (id) of the Fabric. Use the id attribute of the [hyperfabric_fabric](https://registry.terraform.io/providers/cisco-open/hyperfabric/latest/docs/resources/fabric) resource or [hyperfabric_fabric](https://registry.terraform.io/providers/cisco-open/hyperfabric/latest/docs/data-sources/fabric) data source.
* `name` - (string) The name of the VRF.

### Read-Only ###

* `id` - (string) The unique identifier (id) of the VRF in the Fabric.
* `vrf_id` - (string) The unique identifier (id) of the VRF.
* `enabled` - (bool) The enabled state of the VRF.
* `is_default` - (string) The flag that denote if the VRF is the default VRF or not.
* `route_target` - (string) The route target associated with the VRF.
* `metadata` - (map) A map of the Metadata of the VRF:
  * `created_at` - (string) The timestamp when this object was created in [RFC3339](https://datatracker.ietf.org/doc/html/rfc3339#section-5.8) format.
  * `created_by` - (string) The user that created this object.
  * `modified_at` - (string) The timestamp when this object was last modified in [RFC3339](https://datatracker.ietf.org/doc/html/rfc3339#section-5.8) format.
  * `modified_by` - (string) The user that modified this object last.
  * `revision_id` - (string) An integer that represent the current revision of the object.
* `description` - (string) The description is a user defined field to store notes about the VRF.
* `asn` - (integer) The Autonomous System Number (ASN) used for the VRF external connections.
* `vni` - (integer) The VXLAN Network Identifier (VNI) used for the VRF.
* `labels` - (list of strings) A list of user-defined labels that can be used for grouping and filtering objects.
* `annotations` - (list of maps) A list of key-value annotations to store user-defined data including complex data such as JSON.
  * `name` - (string) The name used to uniquely identify the annotation.
  * `value` - (string) The value of the annotation.
  * `data_type` - (string) The type of data stored in the value of the annotation.
      - Possible Values: `STRING`, `INT32`, `UINT32`, `INT64`, `UINT64`, `BOOL`, `TIME`, `UUID`, `DURATION`, `JSON`.

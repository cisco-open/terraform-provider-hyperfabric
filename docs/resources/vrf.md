---
subcategory: "Networking"
layout: "hyperfabric"
page_title: "Nexus Hyperfabric: hyperfabric_vrf"
sidebar_current: "docs-hyperfabric-resource-hyperfabric_vrf"
description: |-
  Manages a VRF in a Nexus Hyperfabric Fabric
---

# hyperfabric_vrf

Manages a VRF in a Nexus Hyperfabric Fabric

A VRF is a virtual-routing-and-forwarding instance that represents a routing table deployed across Nodes in the Fabric to implement VRF-lite and the ability to have multiple routing tables on a single Device. A VRF can be associated to multiple VNIs configured with an Anycast Gateway SVI.

## API Paths ##

* `/fabrics/{fabricId|fabricName}/vrfs` `POST`
* `/fabrics/{fabricId|fabricName}/vrfs/{vrfId|name}` `GET, PUT, DELETE`

## GUI Information ##

* Location: `> Fabrics > {fabric} > Logical network > Route tables (VRF)`

## Example Usage ##

The configuration snippet below creates a VRF with only the required attributes.

```hcl
resource "hyperfabric_vrf" "example_vrf" {
  fabric_id = hyperfabric_fabric.example_fabric.id
  name      = "my-example-vrf"
}
```
The configuration snippet below shows all possible attributes of a VRF.

```hcl
resource "hyperfabric_vrf" "full_example_vrf" {
  fabric_id   = hyperfabric_fabric.example_fabric.id
  name        = "my-full-example-vrf"
  description = "This VRF is part of a Cisco Nexus Hyperfabric"
  asn         = 65002
  vni         = 170
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
* `name` - (string) The name of the VRF.

### Optional ###

* `description` - (string) The description is a user defined field to store notes about the VRF.
* `asn` - (integer) The Autonomous System Number (ASN) used for the VRF external connections.
* `vni` - (integer) The VXLAN Network Identifier (VNI) used for the VRF.
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

## Importing

An existing VRF can be [imported](https://www.terraform.io/docs/import/index.html) into this resource using the following command:

```bash
terraform import hyperfabric_vrf.example_vrf {fabricId|fabricName}/vrfs/{vrfId|name}
```

Starting in Terraform version 1.5, an existing VRF can be imported
using [import blocks](https://developer.hashicorp.com/terraform/language/import) via the following configuration:

```hcl
import {
  id = "{fabricId|fabricName}/vrfs/{vrfId|name}"
  to = hyperfabric_vrf.example_vrf
}
```

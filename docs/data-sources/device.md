---
subcategory: "Devices"
layout: "hyperfabric"
page_title: "Nexus Hyperfabric: hyperfabric_device"
sidebar_current: "docs-hyperfabric-data-source-hyperfabric_device"
description: |-
  Data source for a Nexus Hyperfabric Device
---

# hyperfabric_device

Data source for a Nexus Hyperfabric Device

A Device is a physical device such as a Cisco 6000 switch managed by Cisco Nexus Hyperfabric that can be bound to a Node in a Fabric.

## API Paths ##

* `/devices` `GET`

## GUI Information ##

* Location: `> Devices`

## Example Usage ##

```hcl
data "hyperfabric_device" "example_device" {
  serial_number = "TBD00001"
}
```

## Schema ##

### Required One Of ###

* `serial_number` - (string) The serial number of the Device.
* `device_id` - (string) The device identifier of the Device.

### Read-Only ###

* `id` - (string) The unique identifier (id) of the Device.
* `model_name` - (string) The model name of the Device.
* `fabric_id` - (string) The unique identifier of a Fabric.
* `node_id` - (string) The unique identifier of a Node.
* `os_type` - (string) The operating system type of the Device.
* `rack_id` - (string) The unique identifier of a Rack.
* `roles` - (list of strings) A list of roles associated with the Device.
  - Possible Values: `LEAF`, `SPINE`.

---
subcategory: "Blueprint"
layout: "hyperfabric"
page_title: "Nexus Hyperfabric: hyperfabric_node_management_port"
sidebar_current: "docs-hyperfabric-data-source-hyperfabric_node_management_port"
description: |-
  Data source for a Management Port of a Node in a Nexus Hyperfabric Fabric
---

# hyperfabric_node_management_port

Data source for a Management Port of a Node in a Nexus Hyperfabric Fabric

A Management Port is an Out of Band network interface of a Node used to communicate with the Cisco Nexus Hyperfabric Cloud Controller.

## API Paths ##

* `/fabrics/{fabricId|fabricName}/nodes/{nodeId|nodeName}/managementPorts/{nodeManagementPortId|name}` `GET`

## GUI Information ##

* Location: `> Fabrics > {fabric} > Nodes > {node} > Configure > Management port`

## Example Usage ##

```hcl
data "hyperfabric_node_management_port" "example_node_management_port" {
  node_id = hyperfabric_node.example_node.id
}
```

## Schema ##

### Required ###
* `node_id` - (string) The unique identifier (id) of a Node in a Fabric. Use the id attribute of the [hyperfabric_node](https://registry.terraform.io/providers/cisco-open/hyperfabric/latest/docs/resources/node) resource or [hyperfabric_node](https://registry.terraform.io/providers/cisco-open/hyperfabric/latest/docs/data-sources/node) data source.

### Optional ###

* `name` - (string) The name of the Management Port of the Node.
  - Default: `eth0`

### Read-Only ###

* `id` - (string) The unique identifier (id) of the Node in the Fabric.
* `cloud_urls` - (list of strings) A list of Cloud URLs used by a Node.
* `config_origin` - (string) The source of the configuration, either from the cloud or the device.
  - Possible Values: `CONFIG_ORIGIN_CLOUD`, `CONFIG_ORIGIN_DEVICE`.
* `connected_state` - (string) The connected state denoting if the port has ever successfully connected to the service.
  - Possible Values: `CONNECTED_STATE_NOT_CONNECTED`, `CONNECTED_STATE_CONNECTED`.
* `description` - (string) The description is a user defined field to store notes about the Management Port of the Node.
* `dns_addresses` - (list of strings) A list of DNS IP addresses used by a Node.
* `enabled` - (bool) The enabled state of the Management Port of the Node.
* `ipv4_config_type` - (string) Determines if the IPv4 configuration is static or from DHCP.
  - Possible Values: `CONFIG_TYPE_STATIC`, `CONFIG_TYPE_DHCP`.
* `ipv4_address` - (string) The IPv4 address for the Management Port of the Node.
* `ipv4_gateway` - (string) The IPv4 gateway address for the Management Port of the Node.
* `ipv6_config_type` - (string) Determines if the IPv6 configuration is static or from DHCP.
  - Possible Values: `CONFIG_TYPE_STATIC`, `CONFIG_TYPE_DHCP`.
* `ipv6_address` - (string) The IPv6 address for the Management Port of the Node.
* `ipv6_gateway` - (string) The IPv6 gateway address for the Management Port of the Node.
* `metadata` - (map) A map of the Metadata of the Node:
  * `created_at` - (string) The timestamp when this object was created in [RFC3339](https://datatracker.ietf.org/doc/html/rfc3339#section-5.8) format.
  * `created_by` - (string) The user that created this object.
  * `modified_at` - (string) The timestamp when this object was last modified in [RFC3339](https://datatracker.ietf.org/doc/html/rfc3339#section-5.8) format.
  * `modified_by` - (string) The user that modified this object last.
  * `revision_id` - (string) An integer that represent the current revision of the object.
* `ntp_addresses` - (list of strings) A list of NTP Server IP addresses used by a Node.
* `no_proxy` - (list of strings) A list of IP addresses or domain names that should not be proxied.
* `proxy_address` - (string) The URL for a configured HTTPs proxy for the Node.
* `proxy_credential_id` - (string) The unique identifier (id) of the set of credentials for the proxy.
* `proxy_username` - (string) A username to be used to authenticate to the proxy.
* `proxy_password` - (string) A password to be used to authenticate to the proxy.
<!-- * `labels` - (list of strings) A list of user-defined labels that can be used for grouping and filtering objects.
* `annotations` - (list of maps) A list of key-value annotations to store user-defined data including complex data such as JSON.
  * `name` - (string) The name used to uniquely identify the annotation.
  * `value` - (string) The value of the annotation.
  * `data_type` - (string) The type of data stored in the value of the annotation.
      - Possible Values: `STRING`, `INT32`, `UINT32`, `INT64`, `UINT64`, `BOOL`, `TIME`, `UUID`, `DURATION`, `JSON`. -->

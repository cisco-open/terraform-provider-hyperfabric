---
subcategory: "Blueprint"
layout: "hyperfabric"
page_title: "Nexus Hyperfabric: hyperfabric_node_management_port"
sidebar_current: "docs-hyperfabric-resource-hyperfabric_node_management_port"
description: |-
  Manages a Management Port of a Node in a Nexus Hyperfabric Fabric
---

# hyperfabric_node_management_port

Manages a Management Port of a Node in a Nexus Hyperfabric Fabric

A Management Port is an Out of Band network interface of a Node used to communicate with the Cisco Nexus Hyperfabric Cloud Controller.

## API Paths ##

* `/fabrics/{fabricId|fabricName}/nodes/{nodeId|nodeName}/managementPorts` `POST`
* `/fabrics/{fabricId|fabricName}/nodes/{nodeId|nodeName}/managementPorts/{nodeManagementPortId|name}` `GET, PUT, DELETE`

## GUI Information ##

* Location: `> Fabrics > {fabric} > Nodes > {node} > Configure > Management port`

## Example Usage ##

The configuration snippet below creates a Management Port of a Node with only the required attributes.

```hcl
resource "hyperfabric_node_management_port" "example_node_management_port" {
  node_id = hyperfabric_node.example_node.id
}
```
The configuration snippet below shows all possible attributes of a Management Port of a Node.

```hcl
resource "hyperfabric_node_management_port" "full_example_node_management_port" {
  node_id          = hyperfabric_node.example_node.id
  name             = "eth0"
  ipv4_config_type = "CONFIG_TYPE_STATIC"
  ipv4_address     = "10.0.0.3/24"
  ipv4_gateway     = "10.0.0.254"
  ipv6_config_type = "CONFIG_TYPE_STATIC"
  ipv6_address     = "2001::3/64"
  ipv6_gateway     = "2001::254"
  dns_addresses    = ["8.8.8.8", "1.1.1.1"]
  cloud_urls       = ["https://hyperfabric.cisco.com"]
  ntp_addresses    = ["be.pool.ntp.org", "us.pool.ntp.org"]
  no_proxy         = ["10.0.0.1", "server.local"]
  proxy_address    = "http://proxy.mycompany.com:80"
  proxy_username   = "my_proxy_user"
  proxy_password   = "my_super_secret_password2"
}
```

## Schema ##

### Required ###
* `node_id` - (string) The unique identifier (id) of a Node in a Fabric. Use the id attribute of the [hyperfabric_node](https://registry.terraform.io/providers/cisco-open/hyperfabric/latest/docs/resources/node) resource or [hyperfabric_node](https://registry.terraform.io/providers/cisco-open/hyperfabric/latest/docs/data-sources/node) data source.

### Optional ###
  
* `name` - (string) The name of the Management Port of the Node.
  - Default: `eth0`
* `description` - (string) The description is a user defined field to store notes about the Management Port of the Node.
* `cloud_urls` - (list of strings) A list of Cloud URLs used by a Node.
* `ipv4_config_type` - (string) Determines if the IPv4 configuration is static or from DHCP.
  - Default: `CONFIG_TYPE_DHCP`.
  - Valid Values: `CONFIG_TYPE_STATIC`, `CONFIG_TYPE_DHCP`.
* `ipv4_address` - (string) The IPv4 address for the Management Port of the Node.
* `ipv4_gateway` - (string) The IPv4 gateway address for the Management Port of the Node.
* `ipv6_config_type` - (string) Determines if the IPv6 configuration is static or from DHCP.
  - Default: `CONFIG_TYPE_DHCP`.
  - Valid Values: `CONFIG_TYPE_STATIC`, `CONFIG_TYPE_DHCP`.
* `ipv6_address` - (string) The IPv6 address for the Management Port of the Node.
* `ipv6_gateway` - (string) The IPv6 gateway address for the Management Port of the Node.
* `dns_addresses` - (list of strings) A list of DNS IP addresses used by a Node.
* `ntp_addresses` - (list of strings) A list of NTP Server IP addresses used by a Node.
* `no_proxy` - (list of strings) A list of IP addresses or domain names that should not be proxied.
* `proxy_address` - (string) The URL for a configured HTTPs proxy for the Node.
* `proxy_username` - (string) A username to be used to authenticate to the proxy.
* `proxy_password` - (string) A password to be used to authenticate to the proxy.

<!-- * `labels` - (list of strings) A list of user-defined labels that can be used for grouping and filtering objects.
* `annotations` - (list of maps) A list of key-value annotations to store user-defined data including complex data such as JSON.

  #### Required ####

  * `name` - (string) The name used to uniquely identify the annotation.
  * `value` - (string) The value of the annotation.

  #### Optional ####

  * `data_type` - (string) The type of data stored in the value of the annotation.
      - Default: `STRING`
      - Valid Values: `STRING`, `INT32`, `UINT32`, `INT64`, `UINT64`, `BOOL`, `TIME`, `UUID`, `DURATION`, `JSON`. -->

### Read-Only ###

* `id` - (string) The unique identifier (id) of the Node in the Fabric.
* `enabled` - (bool) The enabled state of the Management Port of the Node.
* `config_origin` - (string) The source of the configuration, either from the cloud or the device.
  - Possible Values: `CONFIG_ORIGIN_CLOUD`, `CONFIG_ORIGIN_DEVICE`.
* `connected_state` - (string) The connected state denoting if the port has ever successfully connected to the service.
  - Possible Values: `CONNECTED_STATE_NOT_CONNECTED`, `CONNECTED_STATE_CONNECTED`.
* `proxy_credential_id` - (string) The unique identifier (id) of the set of credentials for the proxy.
* `metadata` - (map) A map of the Metadata of the Node Management Port:
  * `created_at` - (string) The timestamp when this object was created in [RFC3339](https://datatracker.ietf.org/doc/html/rfc3339#section-5.8) format.
  * `created_by` - (string) The user that created this object.
  * `modified_at` - (string) The timestamp when this object was last modified in [RFC3339](https://datatracker.ietf.org/doc/html/rfc3339#section-5.8) format.
  * `modified_by` - (string) The user that modified this object last.
  * `revision_id` - (string) An integer that represent the current revision of the object.

## Importing

An existing Management Port of a Node can be [imported](https://www.terraform.io/docs/import/index.html) into this resource using the following command:

```bash
terraform import hyperfabric_node_management_port.example_node_management_port {fabricId}/nodes/{nodeId}/managementPorts/{id}
```

Starting in Terraform version 1.5, an existing Management Port of a Node can be imported
using [import blocks](https://developer.hashicorp.com/terraform/language/import) via the following configuration:

```hcl
import {
  id = "{fabricId}/nodes/{nodeId}/managementPorts/{id}"
  to = hyperfabric_node_management_port.example_node_management_port
}
```

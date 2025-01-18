resource "hyperfabric_fabric" "fab1" {
  name        = "terraform-fab1"
  description = "My super description"
}

resource "hyperfabric_fabric" "fab2" {
  name        = "terraform-fab2"
  description = "My super description for Fab2"
  address     = "170 West Tasman Dr."
  city        = "San Jose"
  country     = "USA"
  location    = "sj01-1-101-AAA01"
  labels = [
    "sj01-1-101-AAA01",
    "blue",
    "BT"
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

data "hyperfabric_fabric" "fab2" {
  name = hyperfabric_fabric.fab2.name
}

output "datasource_fabric_description" {
  value = data.hyperfabric_fabric.fab2.description
}

resource "hyperfabric_node" "node1" {
  fabric_id  = hyperfabric_fabric.fab1.id
  name       = "fab1-leaf1"
  model_name = "HF6100-32D"
  roles      = ["LEAF"]
}

resource "hyperfabric_node_management_port" "node1" {
  node_id          = hyperfabric_node.node1.id
  ipv4_config_type = "CONFIG_TYPE_STATIC"
  ipv4_address     = "10.0.0.2/24"
  ipv4_gateway     = "10.0.0.254"
  ipv6_config_type = "CONFIG_TYPE_STATIC"
  ipv6_address     = "2001::1/64"
  ipv6_gateway     = "2001::254"
  dns_addresses    = ["8.8.8.8", "1.1.1.1"]
}

resource "hyperfabric_node" "node2" {
  fabric_id   = hyperfabric_fabric.fab1.id
  name        = "fab1-leaf2"
  description = "The 2nd Leaf of this Fabric2"
  model_name  = "HF6100-32D"
  roles       = ["LEAF"]
}


resource "hyperfabric_node" "node3" {
  fabric_id     = hyperfabric_fabric.fab1.id
  name          = "fab1-spine1"
  model_name    = "HF6100-32D"
  roles         = ["SPINE"]
  serial_number = "ABCDF000"
  description   = "The 1st Spine of this Fabric"
  location      = "SJ01-1-101-AAA01"
  labels        = ["blue", "red", "green"]
  annotations = [
    {
      name  = "color"
      value = "red"
    },
    {
      name      = "rack"
      value     = "AAA01"
      data_type = "STRING"
    }
  ]
}

resource "hyperfabric_node_management_port" "node3" {
  node_id          = hyperfabric_node.node3.id
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

resource "hyperfabric_node_port" "node3_eth1_10" {
  node_id = hyperfabric_node.node3.id
	name = "Ethernet1_10"
  roles = ["UNUSED_PORT"]
}

data "hyperfabric_node_management_port" "node3" {
  node_id = hyperfabric_node_management_port.node3.node_id
}

output "datasource_node_management_port_config_type" {
  value = data.hyperfabric_node_management_port.node3.ipv4_config_type
}


data "hyperfabric_node" "node1" {
  fabric_id = hyperfabric_fabric.fab1.id
  name      = hyperfabric_node.node1.name
}

output "datasource_node_position" {
  value = data.hyperfabric_node.node1.position
}

resource "hyperfabric_node" "node4" {
  fabric_id  = hyperfabric_fabric.fab1.id
  name       = "fab1-spine2"
  model_name = "HF6100-32D"
  roles      = ["SPINE"]
}

resource "hyperfabric_node_management_port" "node4" {
  node_id = hyperfabric_node.node4.id
}


resource "hyperfabric_connection" "node1-node3" {
  fabric_id = hyperfabric_fabric.fab1.id
  local = {
    node_id   = hyperfabric_node.node1.node_id
    port_name = "Ethernet1_1"
  }
  remote = {
    node_id   = hyperfabric_node.node3.node_id
    port_name = "Ethernet1_1"
  }
}

resource "hyperfabric_connection" "node2-node3" {
  fabric_id = hyperfabric_fabric.fab1.id
  local = {
    node_id   = hyperfabric_node.node2.node_id
    port_name = "Ethernet1_1"
  }
  remote = {
    node_id   = hyperfabric_node.node3.node_id
    port_name = "Ethernet1_2"
  }
}

resource "hyperfabric_connection" "node1-node4" {
  fabric_id = hyperfabric_fabric.fab1.id
  local = {
    node_id   = hyperfabric_node.node1.node_id
    port_name = "Ethernet1_2"
  }
  remote = {
    node_id   = hyperfabric_node.node4.node_id
    port_name = "Ethernet1_1"
  }
  description  = "Connection between node1 and node4"
  pluggable    = "SFP-10G-AOC7M"
}

resource "hyperfabric_connection" "node2-node4" {
  fabric_id = hyperfabric_fabric.fab1.id
  local = {
    node_id   = hyperfabric_node.node2.node_id
    port_name = "Ethernet1_2"
  }
  remote = {
    node_id   = hyperfabric_node.node4.node_id
    port_name = "Ethernet1_2"
  }
}

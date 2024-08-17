resource "hyperfabric_vrf" "vrf1" {
  fabric_id   = hyperfabric_fabric.fab2.id
  name        = "VRF1"
  description = "My Super New First VRF"
  asn         = 65001
  vni         = 169
}

resource "hyperfabric_vrf" "vrf2" {
  fabric_id   = hyperfabric_fabric.fab1.id
  name        = "VRF2"
  description = "My Super New Second VRF"
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
      data_type = "STRING"
      name  = "rack"
      value = "AAA01"
    }
  ]
}

resource "hyperfabric_vrf" "vrf3" {
  fabric_id = hyperfabric_fabric.fab1.id
  name      = "VRF4"
}

# data "hyperfabric_vrf" "vrf2" {
#   fabric_id = hyperfabric_fabric.fab1.id
#   name      = hyperfabric_vrf.vrf2.name
# }

# output "datasource_vrf_description" {
#     value = data.hyperfabric_vrf.vrf2.description
# }



resource "hyperfabric_vni" "vni1" {
  fabric_id = hyperfabric_fabric.fab1.id
  name      = "VNI1"
}

resource "hyperfabric_vni" "vni2" {
  fabric_id = hyperfabric_fabric.fab1.id
  name      = "VNI2"
  members = [
    {
      node_id   = hyperfabric_node.node2.node_id
      port_name = "Ethernet1_10"
      vlan_id   = 102
    }
  ]
}

resource "hyperfabric_vni" "vni3" {
  fabric_id = hyperfabric_fabric.fab1.id
  name      = "VNI3"
  description = "My Super New Third VNI"
  vni = 100
  labels = [
    "Ohhhh",
    "Ahhh",
    103
  ]
  annotations = [{
    name  = "position"
    value = "LF0068"
  }]
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
      node_id = hyperfabric_node.node2.node_id
      port_name = "*"
      vlan_id   = 103
    },
    {
      node_id   = hyperfabric_node.node2.node_id
      port_name = "Ethernet1_10"
      vlan_id   = 103
    }
  ]
  vrf_id = hyperfabric_vrf.vrf1.vrf_id
}

# data "hyperfabric_vni" "vni3" {
#   fabric_id = hyperfabric_fabric.fab1.id
#   name      = hyperfabric_vni.vni3.name
# }

# output "datasource_vni_description" {
#     value = data.hyperfabric_vni.vni3.description
# }

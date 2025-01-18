resource "hyperfabric_node_loopback" "node1_loopback10" {
  node_id = hyperfabric_node.node1.id
  name = "Loopback10"
  ipv4_address = "10.1.0.1"
}

resource "hyperfabric_node_loopback" "node1_loopback11" {
  node_id      = hyperfabric_node.node1.id
  name         = "Loopback11"
  description  = "Used for BGP peering"
  ipv4_address = "10.1.0.2"
  ipv6_address = "2001:1::2"
  vrf_id       = hyperfabric_vrf.vrf1.vrf_id
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
      name  = "community"
      value = "AAA01"
    }
  ]
}
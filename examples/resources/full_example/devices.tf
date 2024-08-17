data hyperfabric_device "hyperfabric_device" {
  serial_number = "TFAB31204500"
}

output "datasource_device_leaf1_device_id" {
  value = data.hyperfabric_device.hyperfabric_device.device_id
}
# data hyperfabric_device "leaf2" {
#   serial_number = "TFAB97116234"
# }

# data hyperfabric_device "spine1" {
#   serial_number = "TFAB16626855"
# }

# data hyperfabric_device "spine2" {
#   serial_number = "TFAB56169759"
# }

# resource hyperfabric_bind_to_node "leaf1" {
#   node_id = hyperfabric_node.node1.id
#   device_id = data.hyperfabric_device.hyperfabric_device.id
# }

# resource hyperfabric_bind_to_node "leaf2" {
#   node_id = hyperfabric_node.node2.id
#   device_id = data.hyperfabric_device.leaf2.id
# }

# resource hyperfabric_bind_to_node "spine1" {
#   node_id = hyperfabric_node.node3.id
#   device_id = data.hyperfabric_device.spine1.id
# }

# resource hyperfabric_bind_to_node "spine2" {
#   node_id = hyperfabric_node.node4.id
#   device_id = data.hyperfabric_device.spine2.id
# }

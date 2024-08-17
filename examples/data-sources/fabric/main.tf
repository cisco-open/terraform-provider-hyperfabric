data "hyperfabric_fabric" "example_fabric" {
  name = "terraform-example_fabric"
}

output "fabric_id" {
  value = data.hyperfabric_fabric.example_fabric.id
}

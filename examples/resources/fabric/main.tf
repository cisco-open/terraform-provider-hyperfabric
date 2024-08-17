resource "hyperfabric_fabric" "example_fabric" {
  name        = "terraform-example_fabric"
  description = "A Cisco Nexus Hyperfabric-powered Fabric"
}

resource "hyperfabric_fabric" "full_example_fabric" {
  name        = "terraform-full_example_fabric"
  description = "A Cisco Nexus Hyperfabric-powered full example Fabric"
  address     = "170 West Tasman Dr."
  city        = "San Jose"
  country     = "USA"
  location    = "sj01-1-101-AAA01"
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

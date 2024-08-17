terraform {
  required_providers {
    hyperfabric = {
      source = "cisco-open/hyperfabric"
    }
  }
}

provider "hyperfabric" {
  # token = "<MY_HYPERFABRIC_TOKEN>"
}
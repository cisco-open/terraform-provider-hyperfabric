# The Cisco Nexus Hyperfabric Terraform Provider

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.8
- [Go](https://golang.org/doc/install) >= 1.21

## Building The Provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider using the Go `build` command:

```shell
go build
```

## Using the provider

There is two ways to use the provider you just built from source:
* Define a developer override in your `~/.terraformrc` pointing to this Terraform provider directory
* Install the binary as a plugin by moving it to the right location in the example folder

### Defining a developer override

Create a `.terraformrc` file in your home directory (`~/.terraformrc`) with the following content: 
```hcl
provider_installation {
  dev_overrides {
      "registry.terraform.io/cisco-open/hyperfabric" = "<YOUR_TORTUGA_LOCATION>/Tortuga/cloud/terraform"
  }

  # For all other providers, install them directly from their origin provider
  # registries as normal. If you omit this, Terraform will _only_ use
  # the dev_overrides block, and so no other providers will be available.
  direct {}
}
```
If you already have a `~/.terraformrc` file, add the dev_overrides entry above in your existing list of dev_overrides.
You can now go into any Terraform directory and run `terraform plan` such as `examples/resources/full_example` in this directory. You do not need to run `terraform init` for the hyperfabric provider when using the dev_overrides.

### Install the binary as a plugin

Follow the instructions to [install it as a plugin.](https://www.terraform.io/docs/cli/plugins/index.html) After placing it into your plugins directory or the example folder (`examples/resources/full_example`) plugins directory, run `terraform init` to initialize it.

ex.
```hcl
terraform {
  required_providers {
    hyperfabric = {
      source = "cisco-open/hyperfabric"
    }
  }
}

# Configure provider with your Cisco Nexus Hyperfabric token.
provider "hyperfabric" {
  # Use the HYPERFABRIC_TOKEN environment variable to set your bearer token or use:
  # token = "YOUR_HYPERFABRIC_BEARER_TOKEN"
}

resource "hyperfabric_fabric" "example" {
  name = "example-fabric"
}
```

## Developing the Provider

The Cisco Nexus Hyperfabric provider is build using the [terraform-plugin-framework](https://developer.hashicorp.com/terraform/plugin/framework) provider SDK to benefit from Terraform latest capabilities and uses v6 of the Terraform protocol.

### Pre-Requirements

1. Install latest version of [Go](http://www.golang.org)

### New resources and data-sources

* New resources and data-sources are located in the [internal/provider](https://github.com/cisco-open/terraform-provider-hyperfabric/tree/master/internal/provider) directory.

<!-- * The `provider.go`, `resource_*.go`, `resource_*_test.go`, `data_source_*.go`, `data_source_*_test.go` are generated with templates and should not be changed manually. Files that are automatically generated start with `// Code generated by "gen/generator.go"; DO NOT EDIT.` and should not be modified. When a file is not generated correctly, the template of that file must be adjusted.

* Examples used in the documentation are generated automatically and stored in the [examples/resources](https://github.com/CiscoDevNet/terraform-provider-aci/tree/master/examples/resources) and [examples/data-sources](https://github.com/CiscoDevNet/terraform-provider-aci/tree/master/examples/data-sources) directories.

* Documentation for resources and datasources are generated automatically and stored in the [docs](https://github.com/CiscoDevNet/terraform-provider-aci/tree/master/docs) directory.

* There are a few exceptions of static files which need to be changed manually:

  * Files related to `rest_managed`
  * [provider_test.go](https://github.com/CiscoDevNet/terraform-provider-aci/tree/master/internal/provider/provider_test.go)
  * [provider.tf](https://github.com/CiscoDevNet/terraform-provider-aci/tree/master/examples/provider/provider.tf) -->

## Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).
Please see the Go documentation for the most up to date information about using Go modules.

To add a new dependency `github.com/author/dependency` to your Terraform provider:

```shell
go get github.com/author/dependency
go mod tidy
```

Then commit the changes to `go.mod` and `go.sum`.

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

## Release Notes

See the [Changelog](CHANGELOG.md) for full release notes.

## Related Information

For further information and guides, refer to the following:

- [Cisco Hyperfabric DevNet Dev Center](https://developer.cisco.com/hyperfabric)
- [Cisco Hyperfabric DevNet Documentation](https://developer.cisco.com/docs/hyperfabric)

## License Information

This collection is licensed under the [Mozilla Public License Version 2.0](LICENSE)
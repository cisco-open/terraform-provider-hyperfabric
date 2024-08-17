# The Cisco Nexus Hyperfabric Terraform Provider



_This template repository is built on the [Terraform Plugin Framework](https://github.com/hashicorp/terraform-plugin-framework). The template repository built on the [Terraform Plugin SDK](https://github.com/hashicorp/terraform-plugin-sdk) can be found at [terraform-provider-hyperfabric](https://github.com/cisco-open/terraform-provider-hyperfabric). See [Which SDK Should I Use?](https://developer.hashicorp.com/terraform/plugin/framework-benefits) in the Terraform documentation for additional information._

This repository is a *template* for a [Terraform](https://www.terraform.io) provider. It is intended as a starting point for creating Terraform providers, containing:

- A resource and a data source (`internal/provider/`),
- Examples (`examples/`) and generated documentation (`docs/`),
- Miscellaneous meta files.

These files contain boilerplate code that you will need to edit to create your own Terraform provider. Tutorials for creating Terraform providers can be found on the [HashiCorp Developer](https://developer.hashicorp.com/terraform/tutorials/providers-plugin-framework) platform. _Terraform Plugin Framework specific guides are titled accordingly._

Please see the [GitHub template repository documentation](https://help.github.com/en/github/creating-cloning-and-archiving-repositories/creating-a-repository-from-a-template) for how to create a new repository from this template on GitHub.

Once you've written your provider, you'll want to [publish it on the Terraform Registry](https://developer.hashicorp.com/terraform/registry/providers/publishing) so that others can use it.

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.8
- [Go](https://golang.org/doc/install) >= 1.21

## Building The Provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider using the Go `install` command:

```shell
go install
```

## Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).
Please see the Go documentation for the most up to date information about using Go modules.

To add a new dependency `github.com/author/dependency` to your Terraform provider:

```shell
go get github.com/author/dependency
go mod tidy
```

Then commit the changes to `go.mod` and `go.sum`.

## Using the provider

If you are building the provider, follow the instructions to [install it as a plugin.](https://www.terraform.io/docs/cli/plugins/index.html) After placing it into your plugins directory, run `terraform init` to initialize it.

ex.
```hcl
terraform {
  required_providers {
    aci = {
      source = "ciscodevnet/aci"
    }
  }
}

#configure provider with your cisco aci credentials.
provider "aci" {
  # cisco-aci user name
  username = "admin"
  # cisco-aci password
  password = "password"
  # cisco-aci url
  url      = "https://my-cisco-aci.com"
  insecure = true
  proxy_url = "https://proxy_server:proxy_port"
}

resource "aci_tenant" "test-tenant" {
  name        = "test-tenant"
  description = "This tenant is created by terraform"
}

resource "aci_application_profile" "test-app" {
  tenant_dn   = aci_tenant.test-tenant.id
  name        = "test-app"
  description = "This application profile is created by terraform"
}
```
Note : If you are facing the issue of `invalid character '<' looking for beginning of value` while running `terraform apply`, use signature based authentication in that case, or else use `-parallelism=1` with `terraform plan` and `terraform apply` to limit the concurrency to one thread.

```
terraform plan -parallelism=1
terraform apply -parallelism=1
```  


```hcl
  provider "aci" {
      # cisco-aci user name
      username = "admin"
      # private key path
      private_key = "path to private key"
      # Certificate Name
      cert_name = "user-cert"
      # cisco-aci url
      url      = "https://my-cisco-aci.com"
      insecure = true
  }
```

Note: The value of "cert_name" argument must match the name of the certificate object attached to the APIC user (aaaUserCert) used for signature-based authentication


## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

The Hyperfabric provider is built using the [terraform-plugin-framework](https://developer.hashicorp.com/terraform/plugin/framework) to benefit from Terraform latest capabilities and uses v6 of the Terraform protocol.



* New resources and data-sources are located in the [internal/provider](https://github.com/cisco-open/terraform-provider-hyperfabric/tree/master/internal/provider) directory.

* The `provider.go`, `resource_*.go`, `resource_*_test.go`, `data_source_*.go`, `data_source_*_test.go` are generated with templates and should not be changed manually. Files that are automatically generated start with `// Code generated by "gen/generator.go"; DO NOT EDIT.` and should not be modified. When a file is not generated correctly, the template of that file must be adjusted.

* Examples used in the documentation are generated automatically and stored in the [examples/resources](https://github.com/CiscoDevNet/terraform-provider-aci/tree/master/examples/resources) and [examples/data-sources](https://github.com/CiscoDevNet/terraform-provider-aci/tree/master/examples/data-sources) directories.

* Documentation for resources and datasources are generated automatically and stored in the [docs](https://github.com/CiscoDevNet/terraform-provider-aci/tree/master/docs) directory.

* There are a few exceptions of static files which need to be changed manually:

  * Files related to `rest_managed`
  * [provider_test.go](https://github.com/CiscoDevNet/terraform-provider-aci/tree/master/internal/provider/provider_test.go)
  * [provider.tf](https://github.com/CiscoDevNet/terraform-provider-aci/tree/master/examples/provider/provider.tf)

#### Generator Process

New resources and data-sources are generated with [generator.go](https://github.com/CiscoDevNet/terraform-provider-aci/tree/master/examples/gen/generator.go). The generator assumes that the following directories exist in the directory where the generate.go file is located:

- `./definitions`
  * contains the manually created YAML files with the ACI class definitions for overwriting the meta data retrieved from APIC

- `./meta`
  * contains the JSON files with the ACI class metadata which are retrieved from APIC

- `./templates`
  * contains the Go templates used to generate the full provider code

- `./testvars`
  * contains the generated YAML files with the test variables used in the *_test.go and example files

The below steps should be followed for developing `terraform-plugin-framework` resources and data-sources:

1. Create issue and comment that you will be working on the issue.

2. Fork the terraform-provider-aci repository.

3. Clone the forked code to your local machine.

4. Add minimum required class details to [classes.yaml](https://github.com/CiscoDevNet/terraform-provider-aci/tree/master/examples/gen/definitions/classes.yaml).
  * `ui_location`
  * `sub_category`

5. Copy new meta file(s) or replace existing with newer version of the meta file(s) to the [meta](https://github.com/CiscoDevNet/terraform-provider-aci/tree/master/examples/gen/meta) directory in the following format: `<classname>.json`. Assure that all relationship classes are also added to the `meta` directory. The `GEN_HOST` and `GEN_CLASSES` environment variables can be leveraged to retrieve the classes automatically before rendering the code.
  * `GEN_HOST` can be set to a resolvable APIC host, example: `192.158.1.38`. When not provided the devnet documentation will be retrieved.
  * `GEN_CLASSES` should be set to a comma-separated string of classes, example: `fvTenant` or `fvTenant,fvBD`

6. Run `go generate` in the root of the local repository where the `main.go` is located. The following files should be created or updated:
  * `resource_<resource-name>.go` and `resource_<resource-name>_test.go`
  * `data_source_<resource-name>.go` and `data_source_<resource-name>_test.go`
  * `<resource-name>.md` in the [docs/resources](https://github.com/CiscoDevNet/terraform-provider-aci/tree/master/docs/resources) and [docs/data-sources](https://github.com/CiscoDevNet/terraform-provider-aci/tree/master/docs/data-sources) directories.
  * `resource.tf` and `data-source.tf`
  * `provider.go`

6. Validate generated files generated on the following items:
  * resource name
    * incorrect resource names should be overwritten in [classes.yaml](https://github.com/CiscoDevNet/terraform-provider-aci/tree/master/examples/gen/definitions/classes.yaml) by providing `resource_name: "<resource-name>"`
  * schema attributes ( class properties and children )
    * incorrect property names should be overwritten in [properties.yaml](https://github.com/CiscoDevNet/terraform-provider-aci/tree/master/examples/gen/definitions/properties.yaml) by providing a global or class specific overwrite
    * all relationships classes are included
  * documentation contains the example and all information is correct
  * Run `go generate` again until the files are generated as preferred. When you are not achieving the desired outputs, reource, data-source, example and/or documentation templates should be adjusted.

8. Test the generated code of your resources and data-sources
  * Execute the tests for all your generated resources and data-sources
  * Set the below environment variables when running tests:
    ```sh
    export ACI_VAL_REL_DN=false
    export TF_ACC=1
    export ACI_USERNAME="USER"
    export ACI_URL="https://IPADDRESS"
    export ACI_PASSWORD="PASSWORD"
    ```
  * The following command can be used `go test internal/provider/* -v -run <test-name>`, where the test name can be found in the `resource_<resource-name>_test.go` and `data_source_<resource-name>_test.go` files.
  * Adjust generated test_values in case incorrect values are used. All test values are resolved by estimating the values. Some values are hard to resolve with and thus need to be overwritten in [properties.yaml](https://github.com/CiscoDevNet/terraform-provider-aci/tree/master/examples/gen/definitions/properties.yaml) file under the class specific test_values.
  * Some tests have a dependency on a specific parent, which is unknown because the resource is not an autogenerated resource yet. This can be set in the [properties.yaml](https://github.com/CiscoDevNet/terraform-provider-aci/tree/master/examples/gen/definitions/properties.yaml) file under the the class specific parents.
  * Run `go generate` again until the files are generated as preferred. When you are not achieving the desired outputs, the test templates should be adjusted.

9. Create PR for the code and request review from active maintainers.

10. Review process

### Troubleshooting 

#### Go Generate Fail

If you encounter an error message while generating the resources using go generate, a few steps can be followed:

1. Make sure that you're running the latest version of Go

2. Update dependencies and populate the vendor directory: If you're using Go modules, you can update your dependencies and populate your vendor directory by running the following commands in your terminal:

```sh
go mod tidy
go mod vendor
```

The go mod tidy command will clean up unused dependencies and add missing ones. The go mod vendor command will copy all dependencies into the vendor directory.

3. Disable vendor mode (if necessary): If the error still persists, consider disabling vendor mode by setting GOFLAGS="-mod=mod", and then run go get again:

```sh
export GOFLAGS="-mod=mod"
```

This will force Go to fetch the package directly, regardless of the vendor directory. After the package is downloaded, you can switch back to vendor mode if needed.

#### Missing packages

If you encounter an error indicating that the golang.org/x/text/language package is missing from your vendor directory, you can fetch it by following these steps:

1. Disable vendor mode: Go operates in vendor mode when the -mod=vendor flag is set. You'll need to disable vendor mode to fetch packages directly. Run the following command in your terminal:

```sh
export GOFLAGS="-mod=mod"
```

2. Fetch the missing package: Now that vendor mode is disabled, you can fetch the missing package by running:

```sh
go get golang.org/x/text/language
```

This command tells Go to fetch the golang.org/x/text/language package directly, regardless of the vendor directory

3. Re-enable vendor mode (if necessary): If you wish to switch back to vendor mode, you can do so by running:

```sh
export GOFLAGS="-mod=vendor"
```


### Compiling

To compile the provider, run `make build`. This will build the provider with sanity checks present in scripts directory and put the provider binary in `$GOPATH/bin` directory.

<strong>Important: </strong>To successfully use the provider you need to follow these steps:

- Copy or Symlink the provider from the `$GOPATH/bin` to `~/.terraform.d/plugins/terraform.local/CiscoDevNet/aci/<Version>/<architecture>/` for example:
  ```bash
  ln -s ~/go/bin/terraform-provider-aci ~/.terraform.d/plugins/terraform.local/CiscoDevNet/aci/2.3.0/linux_amd64/terraform-provider-aci
  ```

- Edit the Terraform Provider Configuration to use the local provider.

  ```hcl
  terraform {
    required_providers {
      aci = {
        source = "terraform.local/CiscoDevNet/aci"
        version = "2.3.0"
      }
    }
  }
  ```

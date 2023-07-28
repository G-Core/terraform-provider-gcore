Terraform Gcore Provider
------------------------------
- Terraform provider page: https://registry.terraform.io/providers/G-Core/gcore

<img src="https://gcore.com/img/logo.svg" data-src="https://gcore.com/img/logo.svg" alt="Gcore" width="500px" width="500px"> 
====================================================================================

- [![Gitter chat](https://badges.gitter.im/hashicorp-terraform/Lobby.png)](https://gitter.im/hashicorp-terraform/Lobby)
- Mailing list: [Google Groups](http://groups.google.com/group/terraform-tool)

Requirements
------------

-	[Terraform](https://www.terraform.io/downloads.html) 0.13.x
-	[Go](https://golang.org/doc/install) 1.20 (to build the provider plugin)

Upgrade state
-------
To switch state from the deprecated provider (builds <= 0.3.64) please run
```sh
terraform state replace-provider registry.terraform.io/g-core/gcorelabs registry.terraform.io/g-core/gcore
```

Building the provider
---------------------
```sh
GOPATH=$(go env GOPATH)
mkdir -p $GOPATH/src/github.com/terraform-providers
cd $GOPATH/src/github.com/terraform-providers
git clone https://github.com/G-Core/terraform-provider-gcore.git
cd $GOPATH/src/github.com/terraform-providers/terraform-provider-gcore
make build
```

### Override Terraform provider

To override terraform provider for development goals you do next steps:

- create Terraform configuration file and point provider to development path
- comment out the override if you want to use published binary from the upstream
```shell
cat > ~/.terraformrc << EOF
provider_installation {

  dev_overrides {
      "local.gcore.com/repo/gcore" = "$(go env GOPATH)/src/github.com/terraform-providers/terraform-provider-gcore/bin"
      "registry.terraform.io/g-core/gcore" = "$(go env GOPATH)/src/github.com/terraform-providers/terraform-provider-gcore/bin"
  }

  # For all other providers, install them directly from their origin provider
  # registries as normal. If you omit this, Terraform will _only_ use
  # the dev_overrides block, and so no other providers will be available.
  direct {}
}
EOF
```

Optionally specify `local.gcore.com/repo/gcore` in main.tf configuration file
```shell
terraform {
  required_providers {
    gcore = {
      source  = "G-Core/gcore"
      version = ">= 0.3.65"
      # source = "local.gcore.com/repo/gcore"
      # version = ">=0.3.64"
    }
  }
  required_version = ">= 0.13.0"
}
```

Using the provider
------------------
To use the provider, prepare configuration files based on examples

```sh
$ cp ./examples/... .
$ terraform init # not needed when override is in use
```

Updating docs
-------------
Don't forget to add docs and examples to support your contribution. Update [tfplugindocs](//github.com/hashicorp/terraform-plugin-docs/releases/) when needed.
```sh
$ tfplugindocs
$ git add .
```
Thank You

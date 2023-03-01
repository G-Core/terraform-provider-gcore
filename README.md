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
-	[Go](https://golang.org/doc/install) 1.14 (to build the provider plugin)

Building the provider
---------------------
```sh
$ mkdir -p $GOPATH/src/github.com/terraform-providers
$ cd $GOPATH/src/github.com/terraform-providers
$ git clone https://github.com/G-Core/terraform-provider-gcore.git
$ cd $GOPATH/src/github.com/terraform-providers/terraform-provider-gcore
$ make build
```

### Override Terraform provider

To override terraform provider for development goals you do next steps: 

create Terraform configuration file
```shell
$ touch ~/.terraformrc
```

point provider to development path
```shell
provider_installation { 
 
  dev_overrides { 
      "local.gcore.com/repo/gcore" = "/<dev-path>/terraform-provider-gcore/bin" 
  } 
 
  # For all other providers, install them directly from their origin provider 
  # registries as normal. If you omit this, Terraform will _only_ use 
  # the dev_overrides block, and so no other providers will be available. 
  direct {} 
}
```

add `local.gcore.com/repo/gcore` to .tf configuration file
```shell
terraform {
  required_version = ">= 0.13.0"

  required_providers {
    gcore = {
      source = "local.gcore.com/repo/gcore"
    }
  }
}
```

Using the provider
------------------
To use the provider, prepare configuration files based on examples

```sh
$ cp ./examples/... .
$ terraform init
```

Thank You

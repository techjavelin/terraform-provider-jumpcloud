# JumpCloud Terraform Provider
[![continuous // main](https://github.com/techjavelin/terraform-provider-jumpcloud/actions/workflows/continuous.yml/badge.svg)](https://github.com/techjavelin/terraform-provider-jumpcloud/actions/workflows/continuous.yml)

The JumpCloud Terraform Provider is an unofficial plugin for managing your JumpCloud tenant configuration through the [Terraform](https://www.terraform.io) tool. 

---

## 📚 Documentation

* [Provider - jumpcloud](docs/index.md)
* [Resource - jumpcloud_ad](docs/resources/ad.md)
* [Resource - jumpcloud_devicegroup](docs/resources/devicegroup.md)

### Requirements

* [Terraform](https://terraform.io)
* A [JumpCloud](https://jumpcloud.com) account

### Installation

Terraform uses the [Terraform Registry](https://registry.terraform.io) to download and install providers. To install thisprovider copy and paste the following code into your Terraform configuration

```
terraform {
    required_providers {
        jumpcloud = {
            source = "techjavelin/jumpcloud
            version = ">=1.0.0"
        }
    }
}
```

Then at the command line, run the following command

```
$ terraform init
```

### Configuration

Use of the JumpCloud Provider requires a JumpCloud API Key

#### Getting your API Key
1. As an Administrator or Command Runner, login to the [JumpCloud Console](https://console.jumpcloud.com)
2. From any tag inside of the Admin Console, click your account profile icon in the top-right and select `My API Key` from the drop-down. 
3. Copy the API Key and save it someplace safe. 

Now that you've got your API key, it's time to configure the provider. It is recommended that you use a sensitive variable in your JumpCloud configuration to access the key and provide the value at runtime, so it is never hard-coded into your source code. 

Add the following to your main terraform configuration file (usually `main.tf`)

```
var "jumpcloud_api_key" {
    description = "API Key to access JumpCloud v1, v2, and insights APIs"
    sensitive = true
}
```

You'll also want to update your provider configuration - this can be done at the main level at the module level if your terraform configuration is broken into modules

```
provider "jumpcloud" {
    jumpcloud_api_key = var.jumpcloud_api_key
}
```

To inject the value of the API key at runtime, simple run terraform with the value on the environment 

```
$ TF_VAR_jumpcloud_api_key="<your api key>" terraform <command> [options]
```

*Note, you do not need to provide your key for `init`, `fmt`, or `validate` commands, `plan` and `apply` both require it.*

For example:
```
$ terraform init
$ terraform fmt
$ terraform validate
$ TF_VAR_jumpcloud_api_key="1234" terraform plan --out apply.tfplan
$ TF_VAR_jumpcloud_api_key="1234" terraform apply apply.tfplan
```
#### Rotating your API Key

Occasionally, you may want or need to rotate your API Key. Usually this is due to events such as someone who had access to the value of the API key moving on to a new job or being terminated, simple click the button in the dialog you went to above and update your local storage to reflect the new API key

## 👋 Contributing

🙇 Support & Feedback

### Raise an Issue
If you have found a bug or if you have a feature request, please raise an issue on our issue tracker.

### Vulnerability Reporting
Please do not report security vulnerabilities on the public GitHub issue tracker. Please report directly to oss@techjavelin.com
terraform {
    required_version = "1.9.0"
    required_providers {
        hrobot = {
            source = "registry.terraform.io/suyash-811/hrobot"
        }
    }
}

provider "hrobot" {
    username = ""
    password = ""
}

data "hrobot_vswitch" "example" {
    id = 12345
}

output "example_vswitch" {
  value = data.hrobot_vswitch.example
}
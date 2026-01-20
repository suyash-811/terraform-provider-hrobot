terraform {
  required_providers {
    hrobot = {
      source  = "registry.terraform.io/suyash-811/hrobot"
      version = ">= 0.1.0"
    }
  }
}

provider "hrobot" {
  username = "webservice_username"
  password = "supersecretpassword"
}

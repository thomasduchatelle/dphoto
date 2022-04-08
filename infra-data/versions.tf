terraform {
  required_version = "v1.1.8" # keep it in sync with .terraform-version and terraform cloud.

  required_providers {
    aws = {
      source = "hashicorp/aws"
      version =  "~> 3.75.1"
    }
  }
}
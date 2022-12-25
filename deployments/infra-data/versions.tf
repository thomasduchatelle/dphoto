terraform {
  required_version = "v1.2.5" # keep it in sync with .terraform-version and terraform cloud.

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.22"
    }
  }
}
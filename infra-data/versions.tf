terraform {
  required_version = "v0.15.4"

  required_providers {
    aws = {
      source = "hashicorp/aws"
      version =  "~> 3.74.3"
    }
  }
}
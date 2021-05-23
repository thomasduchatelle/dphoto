terraform {
  backend "remote" {
    organization = "dphoto"

    workspaces {
      prefix = "dphoto-"
    }
  }
}

provider "aws" {
  region = var.region
  assume_role {
    role_arn = "arn:aws:iam::472045025110:role/TerraformRunner"
  }
}
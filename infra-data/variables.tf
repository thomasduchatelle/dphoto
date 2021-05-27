variable "environment_name" {
  description = "Name after which resources will be named"
  type        = string
}

variable "region" {
  description = "AWS Region"
  default     = "eu-west-1"
}

variable "keybase_user" {
  default = "keybase:thomasduchatelle"
}
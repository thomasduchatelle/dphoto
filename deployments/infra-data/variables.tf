variable "environment_name" {
  description = "Name after which resources will be named"
  type        = string
}

variable "keybase_user" {
  default = "keybase:thomasduchatelle"
}

variable "cli_access_keys" {
  description = "List of dates use to roll the access keys"
  default     = ["2022-12"]
  type        = set(string)
}

variable "simple_s3" {
  description = "disable KMS encryption (1$/month), versioning, and glacier retention from S3 store"
  default     = false
}

variable "region" {
  description = "AWS Region"
  default     = "eu-west-1"
}


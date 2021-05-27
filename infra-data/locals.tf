locals {
  path                = "/dphoto/"
  prefix                = "dphoto-${var.environment_name}"
  dynamodb_table_name = "${local.prefix}-index"
  dynamodb_table_arn  = "arn:aws:dynamodb:*:*:table/${local.dynamodb_table_name}"
  tags                = {
    Application = "dphoto-data"
    Environment = var.environment_name
  }
}
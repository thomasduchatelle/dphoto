resource "aws_ssm_parameter" "iam_policy_archive_sns_publish" {
  name  = "/dphoto/${var.environment_name}/iam/policies/archive_sns_publish/arn"
  type  = "String"
  value = aws_iam_policy.archive_sns_publish.arn
  tags  = local.tags
}

resource "aws_ssm_parameter" "iam_policy_archive_sqs_send" {
  name  = "/dphoto/${var.environment_name}/iam/policies/archive_sqs_send/arn"
  type  = "String"
  value = aws_iam_policy.archive_sqs_send.arn
  tags  = local.tags
}

resource "aws_ssm_parameter" "iam_policy_bucket_ro" {
  name  = "/dphoto/${var.environment_name}/iam/policies/storageROArn"
  type  = "String"
  value = aws_iam_policy.storage_ro.arn
  tags  = local.tags
}

resource "aws_ssm_parameter" "iam_policy_dyn_table_rw" {
  name  = "/dphoto/${var.environment_name}/iam/policies/indexRWArn"
  type  = "String"
  value = aws_iam_policy.index_rw.arn
  tags  = local.tags
}

resource "aws_ssm_parameter" "iam_policy_cache_rw" {
  name  = "/dphoto/${var.environment_name}/iam/policies/cacheRWArn"
  type  = "String"
  value = aws_iam_policy.cache_rw.arn
  tags  = local.tags
}

resource "aws_ssm_parameter" "storage_bucket_name" {
  name  = "/dphoto/${var.environment_name}/s3/storage/bucketName"
  type  = "String"
  value = aws_s3_bucket.storage.bucket
  tags  = local.tags
}

resource "aws_ssm_parameter" "cache_bucket_name" {
  name  = "/dphoto/${var.environment_name}/s3/cache/bucketName"
  type  = "String"
  value = aws_s3_bucket.cache.bucket
  tags  = local.tags
}

resource "aws_ssm_parameter" "catalog_table_name" {
  name  = "/dphoto/${var.environment_name}/dynamodb/catalog/tableName"
  type  = "String"
  value = local.dynamodb_table_name
  tags  = local.tags
}

resource "aws_ssm_parameter" "sns_archive_arn" {
  name  = "/dphoto/${var.environment_name}/sns/archive/arn"
  type  = "String"
  value = aws_sns_topic.archive.arn
  tags  = local.tags
}

resource "aws_ssm_parameter" "sqs_archive_arn" {
  name  = "/dphoto/${var.environment_name}/sqs/archive/arn"
  type  = "String"
  value = aws_sqs_queue.async_archive_caching_jobs.arn
  tags  = local.tags
}

resource "aws_ssm_parameter" "sqs_archive_url" {
  name  = "/dphoto/${var.environment_name}/sqs/archive/url"
  type  = "String"
  value = aws_sqs_queue.async_archive_caching_jobs.url
  tags  = local.tags
}

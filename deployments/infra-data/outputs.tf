output "archive_bucket_name" {
  description = "Name of the bucket where medias can be uploaded"
  value       = aws_s3_bucket.storage.bucket
}

output "cache_bucket_name" {
  description = "Name of the bucket where miniatures are cached"
  value       = aws_s3_bucket.cache.bucket
}

output "delegate_access_key_id" {
  description = "AWS access Key to authenticate with dphoto CLI"
  value       = {for k, v in aws_iam_access_key.rolling_cli : k => v.id}
}

output "delegate_secret_access_key" {
  description = "AWS secret access Key to authenticate with dphoto CLI"
  value       = {for k, v in aws_iam_access_key.rolling_cli : k => v.encrypted_secret}
}

output "delegate_secret_access_key_decrypt_cmd" {
  description = "Command to enter to decrypt 'delegate_secret_access_key'"
  value       = "terraform output -raw delegate_secret_access_key | base64 --decode | keybase pgp decrypt"
}

output "dynamodb_name" {
  description = "Name of the table that need to be created"
  value       = local.dynamodb_table_name
}

output "region" {
  description = "AWS Region (from vars)"
  value       = var.region
}

output "sqs_async_archive_jobs_arn" {
  description = "SQS topic used to subscribe lambdas"
  value       = aws_sqs_queue.async_archive_caching_jobs.arn
}

output "sns_archive_arn" {
  description = "SNS topic ARN where are dispatched asynchronous jobs"
  value       = aws_sns_topic.archive.arn
}

output "sqs_archive_url" {
  description = "SQS topic URL where are de-duplicated messages"
  value       = aws_sqs_queue.async_archive_caching_jobs.url
}
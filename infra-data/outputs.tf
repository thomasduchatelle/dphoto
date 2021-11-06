output "bucket_name" {
  description = "Name of the bucket where medias can be uploaded"
  value       = aws_s3_bucket.storage.bucket
}

output "delegate_access_key_id" {
  description = "AWS access Key to authenticate with dphoto CLI"
  value       = aws_iam_access_key.cli.id
}

output "delegate_secret_access_key" {
  description = "AWS secret access Key to authenticate with dphoto CLI"
  value       = aws_iam_access_key.cli.encrypted_secret
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
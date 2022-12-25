resource "aws_iam_user" "cli" {
  name = "${local.prefix}-cli"
  path = local.path
  tags = local.tags
}

resource "aws_iam_user_policy_attachment" "cli_bucket" {
  policy_arn = aws_iam_policy.storage_rw.arn
  user       = aws_iam_user.cli.name
}

resource "aws_iam_user_policy_attachment" "cli_cache" {
  policy_arn = aws_iam_policy.cache_rw.arn
  user       = aws_iam_user.cli.name
}

resource "aws_iam_user_policy_attachment" "cli_table" {
  policy_arn = aws_iam_policy.index_rw.arn
  user       = aws_iam_user.cli.name
}

resource "aws_iam_user_policy_attachment" "archive_sns_publish" {
  policy_arn = aws_iam_policy.archive_sns_publish.arn
  user       = aws_iam_user.cli.name
}

resource "aws_iam_user_policy_attachment" "archive_sqs_send" {
  policy_arn = aws_iam_policy.archive_sqs_send.arn
  user       = aws_iam_user.cli.name
}

resource "aws_iam_access_key" "rolling_cli" {
  for_each = var.cli_access_keys
  user     = aws_iam_user.cli.name
  pgp_key  = var.keybase_user
}
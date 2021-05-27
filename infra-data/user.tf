resource "aws_iam_user" "cli" {
  name = "${local.prefix}-cli"
  path = local.path
  tags = local.tags
}

resource "aws_iam_user_policy_attachment" "cli_bucket" {
  policy_arn = aws_iam_policy.storage_rw.arn
  user       = aws_iam_user.cli.name
}

resource "aws_iam_user_policy_attachment" "cli_table" {
  policy_arn = aws_iam_policy.index_rw.arn
  user       = aws_iam_user.cli.name
}

resource "aws_iam_access_key" "cli" {
  user    = aws_iam_user.cli.name
  pgp_key = var.keybase_user
}
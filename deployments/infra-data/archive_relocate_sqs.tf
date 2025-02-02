resource "aws_sqs_queue" "archive_relocate" {
  name                       = "${local.prefix}-archive-relocate"
  visibility_timeout_seconds = 900*6
  message_retention_seconds  = 14 * 24 * 3600 // 14 days is the maximum
}

resource "aws_iam_policy" "archive_relocate" {
  name   = "${local.prefix}-archive-relocate-sqs-send"
  policy = data.aws_iam_policy_document.archive_relocate.json
  tags   = local.tags
}

data "aws_iam_policy_document" "archive_relocate" {
  statement {
    effect = "Allow"
    actions = [
      "sqs:SendMessage",
    ]
    resources = [
      aws_sqs_queue.archive_relocate.arn
    ]
  }
}
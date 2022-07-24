resource "aws_sns_topic" "archive" {
  name                        = "dphoto-${var.environment_name}-archive-jobs.fifo"
  fifo_topic                  = true
  content_based_deduplication = true
  tags                        = local.tags
}

resource "aws_iam_policy" "archive_sns_publish" {
  name   = "${local.prefix}-archive-sns-publish"
  policy = data.aws_iam_policy_document.archive_sns_publish.json
  tags   = local.tags
}

data "aws_iam_policy_document" "archive_sns_publish" {
  statement {
    effect  = "Allow"
    actions = [
      "sns:Publish",
    ]
    resources = [
      aws_sns_topic.archive.arn
    ]
  }
}

// SQS is required as an intermediate between FIFO SNS and Lambda subscription
resource "aws_sqs_queue" "async_archive_caching_jobs" {
  name                        = "dphoto-${var.environment_name}-async-archive-caching-jobs.fifo"
  fifo_queue                  = true
  content_based_deduplication = true
  visibility_timeout_seconds  = 900 // must be more or equals the function timeout
  tags                        = local.tags
}

resource "aws_sqs_queue_policy" "async_archive_caching_jobs" {
  policy    = data.aws_iam_policy_document.archive_sqs.json
  queue_url = aws_sqs_queue.async_archive_caching_jobs.url
}

data "aws_iam_policy_document" "archive_sqs" {
  statement {
    effect = "Allow"
    principals {
      identifiers = [
        "sns.amazonaws.com",
      ]
      type = "Service"
    }
    actions = [
      "sqs:SendMessage",
    ]
    resources = [
      aws_sqs_queue.async_archive_caching_jobs.arn
    ]
    condition {
      test   = "ArnEquals"
      values = [
        aws_sns_topic.archive.arn,
      ]
      variable = "aws:SourceArn"
    }
  }
}

resource "aws_sns_topic_subscription" "archive_sns_to_sqs" {
  endpoint  = aws_sqs_queue.async_archive_caching_jobs.arn
  protocol  = "sqs"
  topic_arn = aws_sns_topic.archive.arn
}
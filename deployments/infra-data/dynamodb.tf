resource "aws_iam_policy" "index_rw" {
  name   = "${local.prefix}-index-rw"
  path   = local.path
  policy = data.aws_iam_policy_document.index_rw.json
}

data "aws_iam_policy_document" "index_rw" {
  statement {
    effect    = "Allow"
    actions   = [
      "dynamodb:List*",
      "dynamodb:DescribeReservedCapacity*",
      "dynamodb:DescribeLimits",
      "dynamodb:DescribeTimeToLive",
    ]
    resources = [
      "*",
    ]
  }
  statement {
    effect    = "Allow"
    actions   = [
      "dynamodb:BatchGet*",
      "dynamodb:DescribeStream",
      "dynamodb:DescribeTable",
      "dynamodb:Get*",
      "dynamodb:Query",
      "dynamodb:Scan",
      "dynamodb:BatchWrite*",
      "dynamodb:CreateTable",
      "dynamodb:Delete*",
      "dynamodb:Update*",
      "dynamodb:PutItem",
      "dynamodb:TagResource",
    ]
    resources = [
      local.dynamodb_table_arn,
      "${local.dynamodb_table_arn}/*",
    ]
  }
}
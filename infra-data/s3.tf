resource "aws_kms_key" "storage" {
  count = var.simple_s3 ? 0 : 1
  deletion_window_in_days = 30
  tags                    = merge(local.tags, {
    Name = "${local.prefix}-encryption-key"
  })
}

resource "aws_s3_bucket" "storage" {
  bucket = "${local.prefix}-storage"
  acl    = "private"

  dynamic "server_side_encryption_configuration" {
    for_each = var.simple_s3 ? [] : [1]
    content {
      rule {
        apply_server_side_encryption_by_default {
          kms_master_key_id = aws_kms_key.storage.0.arn
          sse_algorithm     = "aws:kms"
        }
      }
    }
  }

  versioning {
    enabled = !var.simple_s3
  }

  lifecycle_rule {
    enabled = !var.simple_s3

    noncurrent_version_transition {
      days          = 0
      storage_class = "GLACIER"
    }

    noncurrent_version_expiration {
      days = 30
    }
  }

  tags = local.tags
}

# Ensure bucket and objects are not public
resource "aws_s3_bucket_public_access_block" "s3_block_public_access" {
  bucket                  = aws_s3_bucket.storage.id
  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

resource "aws_iam_policy" "storage_rw" {
  name   = "${local.prefix}-storage-rw"
  path   = local.path
  policy = data.aws_iam_policy_document.storage_rw.json
}

data "aws_iam_policy_document" "storage_rw" {
  statement {
    effect    = "Allow"
    actions   = [
      "s3:ListBucket",
    ]
    resources = [
      aws_s3_bucket.storage.arn,
    ]
  }
  statement {
    effect    = "Allow"
    actions   = [
      "s3:*Object",
    ]
    resources = [
      "${aws_s3_bucket.storage.arn}/*",
    ]
  }

  dynamic "statement" {
    for_each = var.simple_s3 ? [] : [1]
    content {
      effect    = "Allow"
      actions   = [
        "kms:Decrypt",
        "kms:GenerateDataKey"
      ]
      resources = [
        aws_kms_key.storage.0.arn
      ]
    }
  }
}

resource "aws_iam_policy" "storage_ro" {
  name   = "${local.prefix}-storage-ro"
  path   = local.path
  policy = data.aws_iam_policy_document.storage_ro.json
}

data "aws_iam_policy_document" "storage_ro" {
  statement {
    effect    = "Allow"
    actions   = [
      "s3:ListBucket",
    ]
    resources = [
      aws_s3_bucket.storage.arn,
    ]
  }
  statement {
    effect    = "Allow"
    actions   = [
      "s3:GetObject",
    ]
    resources = [
      "${aws_s3_bucket.storage.arn}/*",
    ]
  }

  dynamic "statement" {
    for_each = var.simple_s3 ? [] : [1]
    content {
      effect    = "Allow"
      actions   = [
        "kms:Decrypt",
        "kms:GenerateDataKey"
      ]
      resources = [
        aws_kms_key.storage.0.arn
      ]
    }
  }
}
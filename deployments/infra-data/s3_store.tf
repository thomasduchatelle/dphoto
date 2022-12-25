resource "aws_kms_key" "storage" {
  count                   = var.simple_s3 ? 0 : 1
  deletion_window_in_days = 30
  tags                    = merge(local.tags, {
    Name = "${local.prefix}-encryption-key"
  })
}

# Ensure bucket and objects are not public
resource "aws_s3_bucket_public_access_block" "s3_storage_block_public_access" {
  bucket                  = aws_s3_bucket.storage.id
  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

resource "aws_s3_bucket_ownership_controls" "storage" {
  bucket = aws_s3_bucket.storage.id

  rule {
    object_ownership = "BucketOwnerEnforced"
  }
}

resource "aws_s3_bucket" "storage" {
  bucket = "${local.prefix}-storage"
  tags   = local.tags
}

resource "aws_s3_bucket_acl" "storage" {
  bucket = aws_s3_bucket.storage.id
  acl    = "private"
}

resource "aws_s3_bucket_server_side_encryption_configuration" "storage" {
  for_each = toset([for k in aws_kms_key.storage : k.arn])
  bucket   = aws_s3_bucket.storage.id
  rule {
    apply_server_side_encryption_by_default {
      kms_master_key_id = each.key
      sse_algorithm     = "aws:kms"
    }
  }
}

resource "aws_s3_bucket_versioning" "storage" {
  bucket = aws_s3_bucket.storage.id
  versioning_configuration {
    status = var.simple_s3 ? "Suspended" : "Enabled"
  }
}

resource "aws_s3_bucket_lifecycle_configuration" "storage" {
  bucket = aws_s3_bucket.storage.id
  rule {
    id     = "deleted-eviction"
    status = "Enabled"

    noncurrent_version_transition {
      noncurrent_days = 0
      storage_class   = "GLACIER"
    }

    noncurrent_version_expiration {
      noncurrent_days = 30
    }
  }

  rule {
    id     = "current-cost-saving"
    status = var.simple_s3 ? "Disabled" : "Enabled"

    transition {
      days          = 7
      storage_class = "GLACIER_IR"
    }
  }
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
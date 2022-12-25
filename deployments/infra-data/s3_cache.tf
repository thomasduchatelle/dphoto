resource "aws_s3_bucket" "cache" {
  bucket = "${local.prefix}-cache"
  tags   = local.tags
}

resource "aws_s3_bucket_acl" "cache" {
  bucket = aws_s3_bucket.cache.id
  acl    = "private"
}

# Ensure bucket and objects are not public
resource "aws_s3_bucket_public_access_block" "s3_cache_block_public_access" {
  bucket                  = aws_s3_bucket.cache.id
  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

resource "aws_s3_bucket_server_side_encryption_configuration" "cache" {
  for_each = toset([for k in aws_kms_key.storage : k.arn])
  bucket   = aws_s3_bucket.cache.id
  rule {
    apply_server_side_encryption_by_default {
      kms_master_key_id = each.key
      sse_algorithm     = "aws:kms"
    }
  }
}

resource "aws_s3_bucket_ownership_controls" "cache" {
  bucket = aws_s3_bucket.cache.id

  rule {
    object_ownership = "BucketOwnerEnforced"
  }
}

resource "aws_s3_bucket_versioning" "cache" {
  bucket = aws_s3_bucket.cache.id
  versioning_configuration {
    status = "Disabled"
  }
}

resource "aws_s3_bucket_lifecycle_configuration" "cache_expiration" {
  bucket = aws_s3_bucket.cache.id
  rule {
    id     = "1-month-eviction"
    status = "Enabled"
    filter {
      prefix = "w="
    }

    expiration {
      // with: GI retrieval $0.03 / GB, STD storage $0.023 / GB / month, and store / cache = 4 times
      // break even point is => 5 months
      days = 120
    }
  }
}

resource "aws_iam_policy" "cache_rw" {
  name   = "${local.prefix}-cache-rw"
  path   = local.path
  policy = data.aws_iam_policy_document.cache_rw.json
}

data "aws_iam_policy_document" "cache_rw" {
  statement {
    effect  = "Allow"
    actions = [
      "s3:ListBucket",
    ]
    resources = [
      aws_s3_bucket.cache.arn,
    ]
  }
  statement {
    effect  = "Allow"
    actions = [
      "s3:*Object",
    ]
    resources = [
      "${aws_s3_bucket.cache.arn}/*",
    ]
  }
}
